package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(datasource string) {
	var err error
	DB, err = sql.Open("sqlite3", datasource)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTABLESQL := `
	CREATE TABLE IF NOT EXISTS CUSTOM_RUGS (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME TEXT NOT NULL,
		EMAIL TEXT NOT NULL,
		DETAILS TEXT NOT NULL,
		STATUS TEXT NOT NULL DEFAULT 'PENDING',
		CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UPDATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	_, err = DB.Exec(createTABLESQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	createTABLESQL = `
	CREATE TABLE IF NOT EXISTS Users (
		ID UUID PRIMARY KEY,
		NAME TEXT NOT NULL,
		EMAIL TEXT NOT NULL,
		PASS_HASH TEXT NOT NULL,
		CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	_, err = DB.Exec(createTABLESQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Database initialized successfully")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
		log.Println("Database closed successfully")
	}
}
