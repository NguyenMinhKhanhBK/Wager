package utils

import (
	"encoding/json"
	"net/http"
)

// Helper function to response an HTTP request
type HTTPUtils interface {
	ReplyJSON(w http.ResponseWriter, data interface{}, code int)
}

type httpUtils struct{}

func NewHTTPUtils() HTTPUtils {
	return &httpUtils{}
}

func (u *httpUtils) ReplyJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
