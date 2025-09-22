package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("\n\n\n\nThe Infinite Library server!\n\n\n\n")

	db_username := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db_username, db_password, db_host, db_name)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Println("\n\nSomeone from their browser wants to join The Infinite Library!!\n\n")
		fmt.Fprintln(w, "\n\nServer saw that you want to join The Infinite Library! :) \n\n")
		_, err = db.Exec(`insert into til_member (id, username) values (12132, "janeaausten")`)
		if err != nil {
			fmt.Println("Insert failed: ", err)
		}
	})

	http.ListenAndServe(":8080", nil)
}
