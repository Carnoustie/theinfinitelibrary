package main

import (
	"fmt"
	"net/http"

	"github.com/Carnoustie/theinfinitelibrary-backend/handlers"
	"github.com/Carnoustie/theinfinitelibrary-backend/repository"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Print("\n\n\n\nThe Infinite Library Server is running\n\n\n\n")

	//Listen to incoming messages in chatrooms
	go broadcaster(handlers.MainChannel, handlers.ChatRoomChannels)

	repository.InitDatabase()

	//responds to health-checks in the cluster
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/signup", handlers.SignupHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/addbook", handlers.AddBookHandler)
	http.HandleFunc("/api/getbooks", handlers.GetBooksHandler)
	http.HandleFunc("/api/postMessage/", handlers.PostMessageHandler)
	http.HandleFunc("/api/chatRoom/", handlers.ChatRoomHandler)

	http.ListenAndServe(":8000", nil)
}
