package main

import (
	"database/sql"
	"log"
	"net/http"
	"wager/database"
	"wager/handlers"
	"wager/middleware"

	"github.com/gorilla/mux"
)

const (
	SERVER_ADDRESS = "127.0.0.1:8080"
	MYSQL_DRIVER   = "mysql"
)

func main() {
	log.Println("Running HTTP server at", SERVER_ADDRESS)
	startHTTPServer()
}

func initDatabase(dataSourceName string) (*database.DB, error) {
	dbConn, err := sql.Open(MYSQL_DRIVER, dataSourceName)
	if err != nil {
		return nil, err
	}
	db := database.NewDB(dbConn)
	return db, nil
}

func startHTTPServer() {
	handler := handlers.NewHandler()

	router := mux.NewRouter()
	router.HandleFunc("/wagers", handler.HandleGetWagers).Methods(http.MethodGet)
	router.HandleFunc("/wagers", handler.HandlePlaceWager).Methods(http.MethodPost)
	router.HandleFunc("/buy/{wager_id}", handler.HandleBuyWager).Methods(http.MethodPost)

	router.Use(middleware.LoggingMiddleware)

	log.Fatal(http.ListenAndServe(SERVER_ADDRESS, router))
}
