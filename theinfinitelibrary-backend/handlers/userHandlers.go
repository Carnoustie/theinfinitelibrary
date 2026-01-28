package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Carnoustie/theinfinitelibrary-backend/repository"
	"golang.org/x/crypto/argon2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //allow sharing response with client

	var bodyContents []byte
	bodyContents, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("\n\nfailed parsing login request with error %s", err)
		return
	}

	var u User
	err = json.Unmarshal(bodyContents, &u) //allow sharing response with client
	if err != nil {
		fmt.Printf("\n\nJSON parsing in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}

	//fetching user from DB
	salt, pwHash, err := repository.GetSaltAndPwHash(u.Username)
	if err != nil {
		fmt.Printf("DB lookup of user in login failed with error: %s", err)
		_, _ = w.Write([]byte("error"))
		return
	} else {
		//Validate password used at login
		encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)
		passwordValidation := subtle.ConstantTimeCompare(pwHash, encryptedPassword) //conceal comparison match times for security
		if passwordValidation == 1 {
			fmt.Printf("\n\nUser %s retrieved from database and password successfully validated for login", u.Username)
			_, _ = w.Write([]byte("Welcome Back! :)"))
		} else {
			fmt.Printf("\n\nUser %s retrieved from database but password was invalid for login", u.Username)
			_, _ = w.Write([]byte("Wrong Password!"))
			return
		}
	}
}



func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //allow sharing response with client
	var bodyContents []byte
	bodyContents, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("\n\nfailed reading membership signup request with error %s", err)
		return
	}
	var u User
	err = json.Unmarshal(bodyContents, &u) //parse json into user
	if err != nil {
		fmt.Printf("\n\nJSON parsing in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}
	//append salt to password to circumvent duplication_
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	//encrypt password for integrity
	encryptedPassword := argon2.IDKey([]byte(u.Password), salt, 1, 64*1024, 4, 32)
	err = repository.AddNewUser(u.Username, salt, encryptedPassword)
	if err != nil {
		fmt.Println("Insert failed: ", err)
		w.Write([]byte("User already exists, pick a different username."))
	} else {
		w.Write([]byte("Welcome to The Infinite Library! :)"))
	}
}
