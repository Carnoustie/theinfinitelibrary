package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
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

	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var bodyContents []byte
		bodyContents, _ = io.ReadAll(r.Body)
		var u User
		_ = json.Unmarshal(bodyContents, &u)

		salt := make([]byte, 16)
		_, _ = rand.Read(salt)
		encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)
		fmt.Println("\n\nencrypted pw: ", encryptedPassword)

		_, err = db.Exec("insert into til_member(username, salt, password_hash) values (?,?,?)", u.Username, salt, encryptedPassword)

		//_, err = db.Exec(`insert into til_member (id, username) values (12132, "janeaausten")`)
		if err != nil {
			fmt.Println("Insert failed: ", err)
			w.Write([]byte("User already exists, pick a different username."))
		} else {
			w.Write([]byte("Welcome to The Infinite Library! :)"))
		}

	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fmt.Println("\n\n\n\nReturning bookworm! :) \n\n")

		var bodyContents []byte
		bodyContents, _ = io.ReadAll(r.Body)
		var u User
		_ = json.Unmarshal(bodyContents, &u)

		var salt []byte
		var pwHash []byte

		err = db.QueryRow("select salt,password_hash from til_member where username=?", u.Username).Scan(&salt, &pwHash)
		if err != nil {
			fmt.Println("DB lookup failed with error: ", err)
			_, _ = w.Write([]byte("Could not find user."))
		} else {

			fmt.Println("\n\nFetched user: ", salt, hex.EncodeToString(pwHash))

			encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)
			fmt.Println("\n\nencrypted pw: ", hex.EncodeToString(encryptedPassword))

			passwordValidation := subtle.ConstantTimeCompare(pwHash, encryptedPassword)

			if passwordValidation == 1 {
				_, _ = w.Write([]byte("Welcome Back! :)"))
			} else {
				_, _ = w.Write([]byte("Wrong Password!"))
			}

			fmt.Println("\n\nBoolean check:\n\n", passwordValidation)
		}

	})

	http.ListenAndServe(":8080", nil)
}
