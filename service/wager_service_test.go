package service

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"wager/conf"
	"wager/mocks"
	"wager/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockSQLResult struct {
	lastInsertedId int64
	err            error
	rowsAffected   int64
}

func (r *mockSQLResult) LastInsertId() (int64, error) {
	return r.lastInsertedId, r.err
}

func (r *mockSQLResult) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

func NewMockWagerService(ctrl *gomock.Controller) (WagerService, *mocks.MockDBManager) {
	mockDb := mocks.NewMockDBManager(ctrl)
	wagerService := &wagerService{
		config: conf.GetDefaultConfig(),
		db:     mockDb,
	}
	return wagerService, mockDb
}

func Test_GetWagerList_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService, _ := NewMockWagerService(ctrl)

	requests := []model.GetWagerListRequest{
		{Page: 0, Limit: 1},
		{Page: 1, Limit: 0},
	}

	for _, req := range requests {
		list, err := mockService.GetWagerList(req)
		assert.Nil(t, list)
		assert.Error(t, err)
	}
}

func Test_GetWagerList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService, mockDB := NewMockWagerService(ctrl)
	req := model.GetWagerListRequest{
		Page:  1,
		Limit: 2,
	}
	mockRows := mocks.NewMockDBRows(ctrl)

	mockDB.EXPECT().Query(gomock.Any(), 2, 0).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true).Times(2)
	mockRows.EXPECT().Scan(gomock.Any()).Times(2)
	mockRows.EXPECT().Next().Return(false)
	mockRows.EXPECT().Close()

	res, err := mockService.GetWagerList(req)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Wagers))
}

func Test_CreateWager_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	wagerService, mockDB := NewMockWagerService(ctrl)
	req := model.CreateWagerRequest{
		TotalWagerValue:   1,
		Odds:              1,
		SellingPercentage: 1,
		SellingPrice:      1,
	}

	mockResult := &mockSQLResult{}

	mockDB.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(mockResult, errors.New("custom error"))
	_, err := wagerService.CreateWager(req)
	assert.Error(t, err)
}

func Test_CreateWager_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	wagerService, mockDB := NewMockWagerService(ctrl)
	req := model.CreateWagerRequest{
		TotalWagerValue:   1,
		Odds:              1,
		SellingPercentage: 1,
		SellingPrice:      1,
	}

	mockResult := &mockSQLResult{
		lastInsertedId: 1,
		err:            nil,
	}

	mockDB.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(mockResult, nil)

	res, err := wagerService.CreateWager(req)

	assert.Equal(t, uint(mockResult.lastInsertedId), res.ID)
	assert.NoError(t, err)
}

func NewDBMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

// TODO Add more BuyWager tests
func Test_BuyWager_BadRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	//db, mockSQL := NewDBMock()
	wagerService, mockDB := NewMockWagerService(ctrl)
	mockRows := mocks.NewMockDBRows(ctrl)

	req := model.BuyWagerRequest{WagerID: 1, BuyingPrice: 1}

	t.Run("WagerID not found", func(t *testing.T) {
		mockDB.EXPECT().Query(gomock.Any(), req.WagerID).Return(mockRows, nil)
		mockRows.EXPECT().Next().Return(false)
		mockRows.EXPECT().Close()
		_, err := wagerService.BuyWager(req)
		assert.Contains(t, err.Error(), "id not found")
	})
}
