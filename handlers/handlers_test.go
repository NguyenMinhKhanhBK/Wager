package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	errorcode "wager/error_code"
	"wager/mocks"
	"wager/model"
	"wager/utils"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

/*
	1. Failed to parse page number
	2. Failed to parse limit numer
	3. No page or limit
	4. Valid page and limit
*/

type MockHandler struct {
	mockWagerService *mocks.MockWagerService
	mockHTTPUtils    *mocks.MockHTTPUtils
}

func NewMockHandler(ctrl *gomock.Controller) (*Handler, *MockHandler) {
	mockHandler := MockHandler{
		mockWagerService: mocks.NewMockWagerService(ctrl),
		mockHTTPUtils:    mocks.NewMockHTTPUtils(ctrl),
	}

	handlers := Handler{
		wagerService: mockHandler.mockWagerService,
		httpUtils:    mockHandler.mockHTTPUtils,
	}

	return &handlers, &mockHandler
}

func Test_HandleGetWagers_BadRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	// invalid page
	req1, err := http.NewRequest("GET", "/wagers?page=a&limit=1", nil)
	assert.NoError(t, err)

	// invalid limit
	req2, err := http.NewRequest("GET", "/wagers?page=1&limit=a", nil)
	assert.NoError(t, err)

	// page is 0
	req3, err := http.NewRequest("GET", "/wagers?page=0&limit=10", nil)
	assert.NoError(t, err)

	// limit is 0
	req4, err := http.NewRequest("GET", "/wagers?page=10&limit=0", nil)
	assert.NoError(t, err)

	// both page and limit are 0
	req5, err := http.NewRequest("GET", "/wagers?page=0&limit=0", nil)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		request       *http.Request
		expectedError errorcode.ErrorResponse
	}{
		{
			name:          "Invalid page number",
			request:       req1,
			expectedError: errorcode.ErrorResponse{Error: "failed to parse page number"},
		},
		{
			name:          "Invalid limit number",
			request:       req2,
			expectedError: errorcode.ErrorResponse{Error: "failed to parse limit number"},
		},
		{
			name:          "Page is 0",
			request:       req3,
			expectedError: errorcode.ErrorResponse{Error: []string{"Page must be larger than 0"}},
		},
		{
			name:          "Limit is 0",
			request:       req4,
			expectedError: errorcode.ErrorResponse{Error: []string{"Limit must be larger than 0"}},
		},
		{
			name:    "Both page and limit are 0",
			request: req5,
			expectedError: errorcode.ErrorResponse{Error: []string{
				"Page must be larger than 0",
				"Limit must be larger than 0",
			}},
		},
	}

	httpHandler := http.HandlerFunc(handler.HandleGetWagers)

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), testcase.expectedError, http.StatusBadRequest)
			httpHandler.ServeHTTP(rr, testcase.request)
		})
	}
}

func Test_HandleGetWagers_DefaultParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	req, err := http.NewRequest("GET", "/wagers", nil)
	assert.NoError(t, err)

	httpHandler := http.HandlerFunc(handler.HandleGetWagers)
	rr := httptest.NewRecorder()

	resp := &model.GetWagerListResponse{Wagers: []model.Wager{
		{ID: 1},
		{ID: 2},
	}}

	mockHandler.mockWagerService.EXPECT().GetWagerList(model.GetWagerListRequest{Page: DEFAULT_PAGE, Limit: DEFAULT_LIMIT}).Return(
		resp,
		nil,
	)
	mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), resp.Wagers, http.StatusOK)
	httpHandler.ServeHTTP(rr, req)
}

func Test_HandleGetWagers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	req, err := http.NewRequest("GET", "/wagers?page=2&limit=20", nil)
	assert.NoError(t, err)

	httpHandler := http.HandlerFunc(handler.HandleGetWagers)
	rr := httptest.NewRecorder()

	resp := &model.GetWagerListResponse{Wagers: []model.Wager{
		{ID: 1},
		{ID: 2},
	}}

	mockHandler.mockWagerService.EXPECT().GetWagerList(model.GetWagerListRequest{Page: 2, Limit: 20}).Return(
		resp,
		nil,
	)
	mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), resp.Wagers, http.StatusOK)
	httpHandler.ServeHTTP(rr, req)
}

func Test_HandlePlaceWager_BadRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	testCases := []struct {
		name          string
		request       model.CreateWagerRequest
		expectedError errorcode.ErrorResponse
	}{
		{
			name:          "Invalid TotalWagerValue and Odds",
			request:       model.CreateWagerRequest{TotalWagerValue: 0, Odds: 0, SellingPercentage: 1, SellingPrice: 1},
			expectedError: errorcode.ErrorResponse{Error: []string{"TotalWagerValue must be larger than 0", "Odds must be larger than 0"}},
		},
		{
			name:          "SellingPrice has more than 2 decimals",
			request:       model.CreateWagerRequest{TotalWagerValue: 1, Odds: 1, SellingPercentage: 1, SellingPrice: 1.111111},
			expectedError: errorcode.ErrorResponse{Error: []string{"SellingPrice must be in monetary format with maximum 2 decimal places"}},
		},
		{
			name:          "SellingPercentage less than 1",
			request:       model.CreateWagerRequest{TotalWagerValue: 1, Odds: 1, SellingPercentage: 0, SellingPrice: 1.11},
			expectedError: errorcode.ErrorResponse{Error: []string{"SellingPercentage must be larger than or equal 1"}},
		},
		{
			name:          "SellingPercentage larger than 100",
			request:       model.CreateWagerRequest{TotalWagerValue: 1, Odds: 1, SellingPercentage: 101, SellingPrice: 1.11},
			expectedError: errorcode.ErrorResponse{Error: []string{"SellingPercentage must be less than or equal 100"}},
		},
		{
			name:          "SellingPrice less than TotalWagerValue * SellingPercentage",
			request:       model.CreateWagerRequest{TotalWagerValue: 5, Odds: 1, SellingPercentage: 100, SellingPrice: 1},
			expectedError: errorcode.ErrorResponse{Error: []string{"SellingPrice must be larger than TotalWagerValue * SellingPercentage"}},
		},
	}

	httpHandler := http.HandlerFunc(handler.HandlePlaceWager)

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			bodyJson, _ := json.Marshal(testcase.request)
			req, err := http.NewRequest(http.MethodPost, "/wagers", bytes.NewReader(bodyJson))
			assert.NoError(t, err)
			mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), testcase.expectedError, http.StatusBadRequest)
			httpHandler.ServeHTTP(rr, req)

		})
	}

}

func Test_HandlePlaceWager_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	httpHandler := http.HandlerFunc(handler.HandlePlaceWager)

	requestBody := model.CreateWagerRequest{TotalWagerValue: 1, Odds: 1, SellingPercentage: 1, SellingPrice: 1}
	bodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/wagers", bytes.NewReader(bodyJson))
	assert.NoError(t, err)

	expectedResp := &model.Wager{
		ID:                  1,
		TotalWagerValue:     1,
		Odds:                1,
		SellingPercentage:   1,
		SellingPrice:        1,
		CurrentSellingPrice: 1,
		PercentageSold:      utils.NullUint{},
		AmountSold:          utils.NullFloat64{},
		PlaceAt:             int64(time.Now().UTC().Unix()),
	}

	rr := httptest.NewRecorder()
	mockHandler.mockWagerService.EXPECT().CreateWager(requestBody).Return(expectedResp, nil)
	mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), expectedResp, http.StatusCreated)
	httpHandler.ServeHTTP(rr, req)
}

func Test_BuyWager_BadRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	testCases := []struct {
		name          string
		request       model.BuyWagerRequest
		expectedError errorcode.ErrorResponse
	}{
		{
			name:          "Invalid WagersID",
			request:       model.BuyWagerRequest{WagerID: 0, BuyingPrice: 1},
			expectedError: errorcode.ErrorResponse{Error: []string{"WagerID must be larger than 0"}},
		},
		{
			name:          "Invalid BuyingPrice",
			request:       model.BuyWagerRequest{WagerID: 1, BuyingPrice: 0},
			expectedError: errorcode.ErrorResponse{Error: []string{"BuyingPrice must be larger than 0"}},
		},
	}

	httpHandler := http.HandlerFunc(handler.HandleBuyWager)

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			bodyJson, _ := json.Marshal(testcase.request)
			url := fmt.Sprintf("/buy/%v", testcase.request.WagerID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyJson))
			assert.NoError(t, err)

			// a hack to set gorilla mux vars
			vars := map[string]string{
				"wager_id": fmt.Sprintf("%d", testcase.request.WagerID),
			}
			req = mux.SetURLVars(req, vars)

			mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), testcase.expectedError, http.StatusBadRequest)
			httpHandler.ServeHTTP(rr, req)

		})
	}
}

func Test_BuyWager_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler, mockHandler := NewMockHandler(ctrl)

	httpHandler := http.HandlerFunc(handler.HandleBuyWager)

	rr := httptest.NewRecorder()
	reqBody := model.BuyWagerRequest{WagerID: 1, BuyingPrice: 1}
	bodyJson, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, "buy/1", bytes.NewReader(bodyJson))
	assert.NoError(t, err)

	expectedResp := &model.Purchase{
		PurchaseID:  1,
		WagerID:     1,
		BuyingPrice: 1,
		BoughtAt:    time.Now().UTC().Unix(),
	}

	// a hack to set gorilla mux vars
	vars := map[string]string{
		"wager_id": fmt.Sprintf("%d", reqBody.WagerID),
	}
	req = mux.SetURLVars(req, vars)

	mockHandler.mockWagerService.EXPECT().BuyWager(reqBody).Return(expectedResp, nil)
	mockHandler.mockHTTPUtils.EXPECT().ReplyJSON(gomock.Any(), expectedResp, http.StatusCreated)
	httpHandler.ServeHTTP(rr, req)
}
