package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleGetWagers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, ok := query["page"]
	if !ok {
		log.Println("get page error")
	}

	limit, ok := query["limit"]
	if !ok {
		log.Println("get limit error")
	}
	log.Println("page:", page, "- limit:", limit)
}

func (h *Handler) HandlePlaceWager(w http.ResponseWriter, r *http.Request) {
	log.Println("place wager")
}

func (h *Handler) HandleBuyWager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wagerId, ok := vars["wager_id"]
	if !ok {
		w.Write([]byte(http.ErrBodyNotAllowed.Error()))
	}

	log.Println("wager id:", wagerId)
}
