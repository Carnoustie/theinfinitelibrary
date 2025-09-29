package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/argon2"
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

	type User struct {
		Username string `json:"username"`
		Password string `json:"Password"`
	}

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var bodyContents []byte
		bodyContents, _ = io.ReadAll(r.Body)
		var u User
		_ = json.Unmarshal(bodyContents, &u)

		salt := make([]byte, 16)
		_, _ = rand.Read(salt)
		encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)
		fmt.Println("\n\nencrypted pw: ", encryptedPassword)

		_, err = db.Exec(`insert into til_member (id, username) values (12132, "janeaausten")`)
		if err != nil {
			fmt.Println("Insert failed: ", err)
		}

	})

	http.ListenAndServe(":8080", nil)
}
