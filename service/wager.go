package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"wager/conf"
	"wager/database"
	"wager/model"
	"wager/utils"

	"github.com/sirupsen/logrus"
)

type WagerService interface {
	CreateWager(request model.CreateWagerRequest) (*model.Wager, error)
	GetWagerList(request model.GetWagerListRequest) (*model.GetWagerListResponse, error)
	BuyWager(request model.BuyWagerRequest) (*model.Purchase, error)
}

type wagerService struct {
	config *conf.Config
	db     database.DBManager
}

func NewWagerService(config *conf.Config, db database.DBManager) WagerService {
	return &wagerService{
		config: config,
		db:     db,
	}
}

func (ws *wagerService) CreateWager(request model.CreateWagerRequest) (*model.Wager, error) {
	wager := model.Wager{
		TotalWagerValue:     request.TotalWagerValue,
		Odds:                request.Odds,
		SellingPercentage:   request.SellingPercentage,
		SellingPrice:        request.SellingPrice,
		CurrentSellingPrice: request.SellingPrice,
		PlaceAt:             time.Now().UTC().Unix(),
	}

	err := ws.insertWager(&wager)
	if err != nil {
		return nil, fmt.Errorf("failed to create wager: %v", err)
	}

	return &wager, nil
}

func (ws *wagerService) insertWager(wager *model.Wager) error {
	insertQuery := fmt.Sprintf("INSERT INTO %v (total_wager_value, odds, selling_percentage, selling_price, current_selling_price, place_at) VALUES (?, ?, ?, ?, ?, ?)", ws.config.SQL.WagerTable)
	res, err := ws.db.Exec(insertQuery, wager.TotalWagerValue, wager.Odds, wager.SellingPercentage, wager.SellingPrice, wager.CurrentSellingPrice, wager.PlaceAt)
	if err != nil {
		return fmt.Errorf("failed to add wager: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get wager id: %v", err)
	}

	wager.ID = uint(id)
	return nil
}

func (ws *wagerService) GetWagerList(request model.GetWagerListRequest) (*model.GetWagerListResponse, error) {
	if request.Page == 0 || request.Limit == 0 {
		return nil, errors.New("invalid request params")
	}
	offset := (request.Page - 1) * request.Limit
	return ws.getWagerList(offset, request.Limit)
}

func (ws *wagerService) getWagerList(offset int, limit int) (*model.GetWagerListResponse, error) {
	result := &model.GetWagerListResponse{}
	query := fmt.Sprintf("SELECT * from %v LIMIT ? OFFSET ?", ws.config.SQL.WagerTable)
	rows, err := ws.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get wagers: %v", err)
	}
	defer rows.Close()

	wagerList := make([]model.Wager, 0)
	for rows.Next() {
		if wager, err := ws.scanSingleWager(rows); err == nil {
			wagerList = append(wagerList, *wager)
		}
	}

	logrus.WithField("wager_list", wagerList).Info("getWagerList")
	result.Wagers = wagerList
	return result, nil
}

func (ws *wagerService) scanSingleWager(rows *sql.Rows) (*model.Wager, error) {
	if rows == nil {
		return nil, errors.New("invalid rows object")
	}

	wager := model.Wager{}
	err := rows.Scan(&wager.ID,
		&wager.TotalWagerValue,
		&wager.Odds,
		&wager.SellingPercentage,
		&wager.SellingPrice,
		&wager.CurrentSellingPrice,
		&wager.PercentageSold,
		&wager.AmountSold,
		&wager.PlaceAt)

	if err != nil {
		logrus.WithError(err).Error("scanSingleWager")
		return nil, err
	}

	return &wager, nil
}

func (ws *wagerService) getWagerByID(id uint) (*model.Wager, error) {
	query := fmt.Sprintf("SELECT * from %v WHERE id=?", ws.config.SQL.WagerTable)
	rows, err := ws.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		wager, err := ws.scanSingleWager(rows)
		if err != nil {
			return nil, err
		}
		return wager, nil
	}

	return nil, errors.New("id not found")
}

func (ws *wagerService) BuyWager(request model.BuyWagerRequest) (*model.Purchase, error) {
	wager, err := ws.getWagerByID(request.WagerID)
	if err != nil {
		return nil, err
	}

	if wager == nil {
		return nil, errors.New("invalid wager object")
	}

	if wager.CurrentSellingPrice < request.BuyingPrice {
		logrus.WithFields(logrus.Fields{
			"current_selling_price": wager.CurrentSellingPrice,
			"buying_price":          request.BuyingPrice,
		}).Info("buying_price must be <= selling_price")
		return nil, errors.New("buying price must be equal or smaller than current selling price")
	}

	wager.CurrentSellingPrice -= request.BuyingPrice
	wager.AmountSold.Float64 += request.BuyingPrice
	wager.AmountSold.Valid = true
	wager.PercentageSold = utils.NewNullUint(uint(wager.AmountSold.Float64 / wager.SellingPrice * 100))

	pur, err := ws.buyWager(&request, wager)
	if err != nil {
		logrus.WithError(err).Error("cannot buy wager")
		return nil, err
	}

	return pur, nil
}

func (ws *wagerService) buyWager(request *model.BuyWagerRequest, wager *model.Wager) (*model.Purchase, error) {
	tx, err := ws.db.BeginTx()
	if err != nil {
		logrus.WithError(err).Error("cannot begin transaction")
		return nil, err
	}

	updateQuery := fmt.Sprintf("UPDATE %v SET current_selling_price=?, percentage_sold=?, amount_sold=? WHERE id=?", ws.config.SQL.WagerTable)
	_, err = ws.db.Exec(updateQuery, wager.CurrentSellingPrice, wager.PercentageSold.Uint, wager.AmountSold.Float64, wager.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	purchase := &model.Purchase{
		WagerID:     request.WagerID,
		BuyingPrice: request.BuyingPrice,
		BoughtAt:    time.Now().UTC().Unix(),
	}
	if err := ws.createPurchase(purchase); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return purchase, nil
}

func (ws *wagerService) createPurchase(purchase *model.Purchase) error {
	query := fmt.Sprintf("INSERT INTO %v (wager_id, buying_price, bought_at) VALUES (?, ?, ?)", ws.config.SQL.PurchaseTable)
	res, err := ws.db.Exec(query, purchase.WagerID, purchase.BuyingPrice, purchase.BoughtAt)
	if err != nil {
		return fmt.Errorf("failed to create purchase: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get purchase id: %v", err)
	}

	purchase.PurchaseID = uint(id)
	return nil
}
