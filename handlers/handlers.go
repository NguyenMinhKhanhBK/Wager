package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	errorcode "wager/error_code"
	"wager/model"
	"wager/service"
	"wager/utils"
	"wager/validator"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_PAGE  = 1
	DEFAULT_LIMIT = 10
)

type Handler struct {
	wagerService service.WagerService
	httpUtils    utils.HTTPUtils
}

func NewHandler(wagerSvrc service.WagerService) *Handler {
	return &Handler{
		wagerService: wagerSvrc,
		httpUtils:    utils.NewHTTPUtils(),
	}
}

func (h *Handler) HandleGetWagers(w http.ResponseWriter, r *http.Request) {
	reqPage := DEFAULT_PAGE
	reqLimit := DEFAULT_LIMIT

	query := r.URL.Query()
	if page, ok := query["page"]; ok {
		num, err := strconv.Atoi(page[0])
		if err != nil {
			jsonErr := errorcode.ErrorResponse{Error: "failed to parse page number"}
			h.httpUtils.ReplyJSON(w, jsonErr, http.StatusBadRequest)
			return
		}
		reqPage = num
	}

	if limit, ok := query["limit"]; ok {
		num, err := strconv.Atoi(limit[0])
		if err != nil {
			jsonErr := errorcode.ErrorResponse{Error: "failed to parse limit number"}
			h.httpUtils.ReplyJSON(w, jsonErr, http.StatusBadRequest)
			return
		}
		reqLimit = num
	}

	req := model.GetWagerListRequest{Page: reqPage, Limit: reqLimit}
	if err := validator.Validate(req); err != nil {
		h.httpUtils.ReplyJSON(w, validator.ErrorMsg(err), http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"page":  reqPage,
		"limit": reqLimit,
	}).Info("RequestQuery")

	wagers, err := h.wagerService.GetWagerList(req)
	if err != nil {
		h.httpUtils.ReplyJSON(w, errorcode.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.httpUtils.ReplyJSON(w, wagers.Wagers, http.StatusOK)
}

func (h *Handler) HandlePlaceWager(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	logrus.WithField("Type", contentType).Info("Content-Type")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read request body")
		jsonErr := errorcode.ErrorResponse{Error: []string{err.Error()}}
		h.httpUtils.ReplyJSON(w, jsonErr, http.StatusBadRequest)
		return
	}

	req := model.CreateWagerRequest{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		logrus.WithError(err).Error("failed to unmarshal request body")
		jsonErr := errorcode.ErrorResponse{Error: []string{err.Error()}}
		h.httpUtils.ReplyJSON(w, jsonErr, http.StatusBadRequest)
		return
	}

	if err := validator.Validate(req); err != nil {
		logrus.WithField("error", validator.ErrorMsg(err)).Info("Validate failed")
		h.httpUtils.ReplyJSON(w, validator.ErrorMsg(err), http.StatusBadRequest)
		return
	}

	if req.SellingPrice <= float64(req.TotalWagerValue*req.SellingPercentage)/100 {
		jsonErr := errorcode.ErrorResponse{Error: []string{"SellingPrice must be larger than TotalWagerValue * SellingPercentage"}}
		h.httpUtils.ReplyJSON(w, jsonErr, http.StatusBadRequest)
		return
	}

	wager, err := h.wagerService.CreateWager(req)
	if err != nil {
		jsonErr := errorcode.ErrorResponse{Error: []string{err.Error()}}
		h.httpUtils.ReplyJSON(w, jsonErr, http.StatusInternalServerError)
		return
	}

	h.httpUtils.ReplyJSON(w, wager, http.StatusCreated)
}

func (h *Handler) HandleBuyWager(w http.ResponseWriter, r *http.Request) {
	req := model.BuyWagerRequest{}
	vars := mux.Vars(r)
	wagerIdStr, ok := vars["wager_id"]
	if !ok {
		h.httpUtils.ReplyJSON(w, errorcode.ErrorResponse{Error: "invalid wager id"}, http.StatusBadRequest)
		return
	}

	wagerId, err := strconv.Atoi(wagerIdStr)
	if err != nil {
		h.httpUtils.ReplyJSON(w, errorcode.ErrorResponse{Error: "failed to parse wager id"}, http.StatusBadRequest)
		return
	}

	req.WagerID = uint(wagerId)

	contentType := r.Header.Get("Content-Type")
	logrus.WithField("Type", contentType).Info("Content-Type")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read request body")
		h.httpUtils.ReplyJSON(w, errorcode.ErrorResponse{Error: "failed to read request body"}, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(data, &req); err != nil {
		h.httpUtils.ReplyJSON(w, errorcode.ErrorResponse{Error: "failed to unmarshal request body"}, http.StatusBadRequest)
		return
	}

	if err := validator.Validate(req); err != nil {
		h.httpUtils.ReplyJSON(w, validator.ErrorMsg(err), http.StatusBadRequest)
		return
	}

	res, err := h.wagerService.BuyWager(req)
	if err != nil {
		h.httpUtils.ReplyJSON(w, errorcode.ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	h.httpUtils.ReplyJSON(w, res, http.StatusCreated)
}
