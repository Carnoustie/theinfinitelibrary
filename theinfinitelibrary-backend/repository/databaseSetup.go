package repository

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func InitDatabase() {
	DB_username := os.Getenv("DB_USER")
	DB_password := os.Getenv("DB_PASSWORD")
	DB_name := os.Getenv("DB_NAME")
	DB_host := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", DB_username, DB_password, DB_host, DB_name) //Data Source Name
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("\n\nCould not open Database due to error %s, terminating server", err)
		os.Exit(1)
	}
}
