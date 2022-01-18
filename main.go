package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"wager/conf"
	"wager/database"
	"wager/handlers"
	"wager/middleware"
	"wager/service"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	MYSQL_DRIVER = "mysql"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{})

	config := conf.GetDefaultConfig()
	if config == nil {
		logrus.Fatal("Failed to load config")
	}

	dsn := fmt.Sprintf("%v:%v@%v", config.SQL.Username, config.SQL.Password, config.SQL.DatabaseAddress)

	/*
		if err := migrateDatabase(dsn); err != nil {
			logrus.WithError(err).Fatal("cannot migrate database")
		}
	*/

	db, err := initDatabase(dsn)
	if err != nil {
		logrus.Fatalf("Failed to init database: %v", err)
	}

	logrus.Info("Initialize database successfully")

	startHTTPServer(config, db)
}

func initDatabase(dataSourceName string) (database.DBManager, error) {
	db, err := sql.Open(MYSQL_DRIVER, dataSourceName)
	if err != nil {
		return nil, err
	}

	return database.NewDB(db), nil
}

/*
func migrateDatabase(dsn string) error {
	db, _ := sql.Open("mysql", dsn)
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"mysql",
		driver,
	)

	m.Steps(2)
}
*/

func startHTTPServer(config *conf.Config, db database.DBManager) {
	if config == nil || db == nil {
		log.Fatal("Invalid intializer objects")
	}

	wagerService := service.NewWagerService(config, db)
	handler := handlers.NewHandler(wagerService)

	router := mux.NewRouter()
	router.HandleFunc(config.Handlers.GetWagerList, handler.HandleGetWagers).Methods(http.MethodGet)
	router.HandleFunc(config.Handlers.CreateWager, handler.HandlePlaceWager).Methods(http.MethodPost)
	router.HandleFunc(config.Handlers.BuyWager, handler.HandleBuyWager).Methods(http.MethodPost)

	router.Use(middleware.LoggingMiddleware)

	logrus.Infof("Running HTTP server at :%v", config.ServerPort)
	serverAddr := fmt.Sprintf(":%v", config.ServerPort)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
