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

type User struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

func main() {
	fmt.Println("\n\n\n\nThe Infinite Library Server is running\n\n\n\n")

	db_username := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db_username, db_password, db_host, db_name) //Data Source Name
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	//responds to health-checks in the cluster
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //allow sharing response with client

		var bodyContents []byte
		bodyContents, err = io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var u User
		_ = json.Unmarshal(bodyContents, &u) //parse json into user

		//append salt to password to circumvent duplication_
		salt := make([]byte, 16)
		_, _ = rand.Read(salt)

		//encrypt password for integrity
		encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)

		_, err = db.Exec("insert into til_member(username, salt, password_hash) values (?,?,?)", u.Username, salt, encryptedPassword)
		if err != nil {
			fmt.Println("Insert failed: ", err)
			w.Write([]byte("User already exists, pick a different username."))
		} else {
			w.Write([]byte("Welcome to The Infinite Library! :)"))
		}

	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //allow sharing response with client

		var bodyContents []byte
		bodyContents, err = io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var u User
		_ = json.Unmarshal(bodyContents, &u) //allow sharing response with client

		//fetching user from DB
		var salt []byte
		var pwHash []byte
		err = db.QueryRow("select salt,password_hash from til_member where username=?", u.Username).Scan(&salt, &pwHash)
		if err != nil {
			fmt.Println("DB lookup failed with error: ", err)
			_, _ = w.Write([]byte("Could not find user."))
		} else {
			fmt.Println("\n\nFetched user: ", salt, hex.EncodeToString(pwHash))

			//Validate password used at login
			encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)
			passwordValidation := subtle.ConstantTimeCompare(pwHash, encryptedPassword) //conceal comparison match times for security

			if passwordValidation == 1 {
				_, _ = w.Write([]byte("Welcome Back! :)"))
			} else {
				_, _ = w.Write([]byte("Wrong Password!"))
			}
		}

	})

	http.ListenAndServe(":8080", nil)
}
