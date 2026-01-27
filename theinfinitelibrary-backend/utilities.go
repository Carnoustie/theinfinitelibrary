package main

import (
	"fmt"
	"github.com/Carnoustie/theinfinitelibrary-backend/handlers"
)

// projectwide error variable to use for logging
var err error

// Listens to incoming messages in chatrooms and directs message to correct chatroom
func broadcaster(mainChan chan handlers.ChatPayLoad, chatRooms map[string][]chan string) {
	var payload handlers.ChatPayLoad
	for {
		payload = <-mainChan
		fmt.Println("\n\nRead following message: ", payload.Message, "\n\n")
		for chatid, channels := range chatRooms {
			fmt.Println("writing in channel ", chatid)
			if payload.ChatRoomID == chatid {
				for _, channel := range channels {
					channel <- payload.Username + ": " + payload.Message
					fmt.Printf("\n\nBroadcasted message%s in chatroom channel %s", payload.Message, chatid)
				}
			}
		}
	}
}

// Computes similarity between booktitles to infer equality of books in light of possible incorrectly spelled book titles by users
func levenshteinDistance(s1 string, s2 string, len_s1 int, len_s2 int) int {
	if len_s1 == 0 {
		return len_s2
	} else if len_s2 == 0 {
		return len_s1
	} else if len_s1 == len_s2 && s1[len_s1-1] == s2[len_s2-1] {
		return levenshteinDistance(s1, s2, len_s1-1, len_s2-1)
	} else {
		return 1 + min(levenshteinDistance(s1, s2, len_s1, len_s2-1), min(levenshteinDistance(s1, s2, len_s1-1, len_s2), levenshteinDistance(s1, s2, len_s1-1, len_s2-1)))
	}
}
