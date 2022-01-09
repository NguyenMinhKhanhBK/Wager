package database

import (
	"database/sql"
	"log"
)

const (
	WAGERS_TABLE   = "wagers"
	PURCHASE_TABLE = "purchase"
)

type DB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *DB {
	if db == nil {
		log.Fatal("nil database instance")
	}
	return &DB{db: db}
}

func (w *DB) InsertWager() {

}
