package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Carnoustie/theinfinitelibrary-backend/algorithms"
	"github.com/Carnoustie/theinfinitelibrary-backend/repository"
)



func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client
	var bodyContents []byte
	bodyContents, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("\n\nParsing payload in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}
	var bk repository.Book
	var u User
	err = json.Unmarshal(bodyContents, &bk)
	if err != nil {
		fmt.Printf("\n\nJSON parsing in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}
	err = json.Unmarshal(bodyContents, &u)
	if err != nil {
		fmt.Printf("\n\nJSON parsing in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}
	userID, err := repository.GetUseridByUsername(u.Username)
	if err != nil {
		fmt.Printf("\n\nDatabase lookup of member %s failed with error %s\n\n", u.Username, err)
		return
	}
	err = repository.PersistBook(userID, bk.Title, bk.Author)
	if err != nil {
		fmt.Printf("\n\nDatabase insert of book %s with title %s by member %s failed with error %s\n\n", bk.Title, bk.Author, u.Username, err)
		return
	} else {
		fmt.Printf("\n\nAdded book with title %s and author %s for user %s", bk.Title, bk.Author, u.Username)
	}
}

func AddBookHandler_Test(t * testing.T) int{
	return 1
}



func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client
	w.Header().Set("Content-Type", "application/json")

	var bodyContents []byte
	bodyContents, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("\n\nFailed parsing HTTP request to %s with error %s", r.URL.Path, err)
	}

	var u User
	err = json.Unmarshal(bodyContents, &u)
	if err != nil {
		fmt.Printf("\n\nJSON parsing in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}

	memberID, err := repository.GetUseridByUsername(u.Username)
	if err != nil {
		fmt.Printf("Database lookup of member in HTTP request to %s failed with error: %s", r.URL.Path, err)
		return
	}

	rows, err := repository.GetBooksByMemberID(memberID)
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: %s", r.URL.Path, err)
		return
	}

	var books []repository.Book
	var allBooks []repository.Book
	var tempBook repository.Book
	allBookRows, err := repository.GetAllBooks()
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: %s", r.URL.Path, err)
		return
	}

	for allBookRows.Next() {
		allBookRows.Scan(&tempBook.Id, &tempBook.MemberID, &tempBook.Title, &tempBook.Author, &tempBook.ChatRoomID)
		allBooks = append(allBooks, tempBook)
	}

	var b repository.Book
	var levenshteinDist int
	// var localUserName string
	countMatchesMap := make(map[int64]int)
	var chatRoomId int64
	var checkChatroomID int64
	for rows.Next() {
		rows.Scan(&b.Id, &b.Title, &b.Author)

		for _, localbook := range allBooks {
			if localbook.MemberID != memberID {
				b_front := strings.Split(b.Title, " ")[0]
				localbook_front := strings.Split(localbook.Title, " ")[0]
				firstLevenshtein := algorithms.LevenshteinDistance(b_front, localbook_front, len(b_front), len(localbook_front))
				if firstLevenshtein < 2 {
					levenshteinDist = algorithms.LevenshteinDistance(strings.ReplaceAll(b.Title, " ", "_"), strings.ReplaceAll(localbook.Title, " ", "_"), len(b.Title), len(localbook.Title))
					fmt.Println("\n\nComputed Levenshtein: ", levenshteinDist)
					if levenshteinDist <= 2 {
						countMatchesMap[localbook.MemberID]++
						if countMatchesMap[localbook.MemberID] >= 2 {
							checkChatroomID, err = repository.GetChatroomIdFromBookId(b.Id)
							if err != nil {
								fmt.Printf("Database lookup of chatroom in HTTP request to %s failed with error: %s", r.URL.Path, err)
							}
							if checkChatroomID == 0 {
								err = repository.AddBookToChatroom(b.Title)

								if err != nil {
									fmt.Printf("Database insertion of book into chatroom in HTTP request to %s failed with error: %s", r.URL.Path, err)
								} else {
									chatRoomId, err = repository.GetChatroomIdFromBookTitle(b.Title)
									if err != nil {
										fmt.Printf("Database lookup of chatroom in HTTP request to %s failed with error: %s", r.URL.Path, err)
									}
									err = repository.SetChatroomIdFromMultipleBooks(b.Id, localbook.Id, chatRoomId)
									if err != nil {
										fmt.Printf("Database update of books in HTTP request to %s failed with error: %s", r.URL.Path, err)
									}
								}
							} else {
								err = repository.SetChatroomIdFromSingleBook(localbook.Id, chatRoomId)
								if err != nil {
									fmt.Printf("Database update of books in HTTP request to %s failed with error: %s", r.URL.Path, err)
								}
							}
						}
					}
				}
			}
		}

	}

	fmt.Printf("\n\nCompleted chatroom assignment updates for books in Database in HTTP request to %s", r.URL.Path)
	rows, err = repository.GetBooksByMemberID(memberID)
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: %s", r.URL.Path, err)
	}

	for rows.Next() {
		rows.Scan(&b.Id, &b.MemberID, &b.Title, &b.Author, &b.ChatRoomID)
		books = append(books, b)
	}
	var matchedUsername string
	for matchedMember, numberOfSharedBooks := range countMatchesMap {
		fmt.Println("\n\nsharedBookCount: ", numberOfSharedBooks)
		fmt.Println("\n\nwith member: ", matchedMember)
		fmt.Println("\n\nlogged in member: ", memberID)
		if numberOfSharedBooks >= 2 {
			matchedUsername, err = repository.GetUsernameByID(matchedMember)
			if err != nil {
				fmt.Printf("Database lookup of username in HTTP request to %s failed with error: %s", r.URL.Path, err)
			}

			fmt.Println("\n\nMatched with member: ", matchedUsername)
		}
	}

	if len(books) == 0 {
		books = []repository.Book{}
	}

	jsonBooks, _ := json.Marshal(&books)
	w.Write([]byte(jsonBooks))
}
