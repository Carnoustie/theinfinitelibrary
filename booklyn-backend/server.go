package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("\n\n\n\nBooklyn server!\n\n\n\n")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Println("\n\nSomeone from their browser wants to join the book club!\n\n")
		fmt.Fprintln(w, "\n\nServer saw that you want to become a Booklyn member :) \n\n")
	})

	http.ListenAndServe(":8080", nil)
}
