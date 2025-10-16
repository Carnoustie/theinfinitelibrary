package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("\n\n\n\nThe Infinite Library Server is running\n\n\n\n")

	//Listen to incoming messages in chatrooms
	go broadcaster(mainChannel, chatRoomChannels)

	initDatabase()

	//responds to health-checks in the cluster
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/signup", SignupHandler)
	http.HandleFunc("/api/login", LoginHandler)
	http.HandleFunc("/api/addbook", AddBookHandler)
	http.HandleFunc("/api/getbooks", GetBooksHandler)
	http.HandleFunc("/api/postMessage/", PostMessageHandler)
	http.HandleFunc("/api/chatRoom/", ChatRoomHandler)

	http.ListenAndServe(":8080", nil)
}
