package main

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func initDatabase() {
	db_username := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db_username, db_password, db_host, db_name) //Data Source Name
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("\n\nCould not open Database due to error %s, terminating server", err)
		os.Exit(1)
	}
}
