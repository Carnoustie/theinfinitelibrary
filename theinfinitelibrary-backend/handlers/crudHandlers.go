package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Carnoustie/theinfinitelibrary-backend/algorithms"
	"github.com/Carnoustie/theinfinitelibrary-backend/repository"
)

type Book struct {
	Id         int64  `json:id`
	MemberID   int64  `json:"member_id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	ChatRoomID int64  `json:"chatroom_id"`
}


func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client
	var bodyContents []byte
	bodyContents, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("\n\nParsing payload in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}
	var bk Book
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
	var userID int
	err = repository.DB.QueryRow("select id from til_member where username=?", u.Username).Scan(&userID)
	if err != nil {
		fmt.Printf("\n\nDatabase lookup of member %s failed with error %s\n\n", u.Username, err)
		return
	}
	_, err = repository.DB.Exec("insert into books (member_id, title, author) values(?,?,?)", userID, bk.Title, bk.Author)
	if err != nil {
		fmt.Printf("\n\nDatabase insert of book %s with title %s by member %s failed with error %s\n\n", bk.Title, bk.Author, u.Username, err)
		return
	} else {
		fmt.Printf("\n\nAdded book with title %s and author %s for user %s", bk.Title, bk.Author, u.Username)
	}
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

	var memberID int64
	err := repository.DB.QueryRow("select id from til_member where username=?", u.Username).Scan(&memberID)
	if err != nil {
		fmt.Printf("Database lookup of member in HTTP request to %s failed with error: ", r.URL.Path, err)
		return
	}

	rows, err := repository.DB.Query("select id, title, author from books where member_id=?", memberID)
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: ", r.URL.Path, err)
		return
	}

	var books []Book
	var allBooks []Book
	var tempBook Book
	allBookRows, err := repository.DB.Query("select id, member_id, title, author from books")
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: ", r.URL.Path, err)
		return
	}

	for allBookRows.Next() {
		allBookRows.Scan(&tempBook.Id, &tempBook.MemberID, &tempBook.Title, &tempBook.Author)
		allBooks = append(allBooks, tempBook)
	}

	var b Book
	var levenshteinDist int
	var localUserName string
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
						err = repository.DB.QueryRow("select username from til_member where id=?", localbook.MemberID).Scan(&localUserName)
						if err != nil {
							fmt.Printf("Database lookup of user %d in HTTP request to %s failed with error: ", localbook.MemberID, r.URL.Path, err)
						}

						countMatchesMap[localbook.MemberID]++
						if countMatchesMap[localbook.MemberID] >= 2 {
							err = repository.DB.QueryRow("select chatroom_id from books where id=?", b.Id).Scan(&checkChatroomID)
							if err != nil {
								fmt.Printf("Database lookup of chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
							}

							// fmt.Println("\n\nfetched: ", checkChatroomID)
							if checkChatroomID == 0 {
								_, err = repository.DB.Exec("insert into chatroom (book_title) values (?)", b.Title)

								if err != nil {
									fmt.Printf("Database insertion of book into chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
								} else {
									err = repository.DB.QueryRow("select id from chatroom where book_title = ?", b.Title).Scan(&chatRoomId)
									if err != nil {
										fmt.Printf("Database lookup of chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
									}

									_, err = repository.DB.Exec("update books set chatroom_id = ? where id = ? or id = ?", chatRoomId, b.Id, localbook.Id)
									if err != nil {
										fmt.Printf("Database update of books in HTTP request to %s failed with error: ", r.URL.Path, err)
									}
								}
							} else {
								_, err = repository.DB.Exec("update books set chatroom_id = ? where id = ?", chatRoomId, localbook.Id)
								if err != nil {
									fmt.Printf("Database update of books in HTTP request to %s failed with error: ", r.URL.Path, err)
								}
							}
						}
					}
				}
			}
		}

	}

	fmt.Println("\n\nCompleted chatroom assignment updates for books in Database in HTTP request to %s", r.URL.Path)
	rows, err = repository.DB.Query("select id, title, author, chatroom_id from books where member_id=?", memberID)
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: ", r.URL.Path, err)
	}

	for rows.Next() {
		rows.Scan(&b.Id, &b.Title, &b.Author, &b.ChatRoomID)
		books = append(books, b)
	}
	var matchedUsername string
	for matchedMember, numberOfSharedBooks := range countMatchesMap {
		fmt.Println("\n\nsharedBookCount: ", numberOfSharedBooks)
		fmt.Println("\n\nwith member: ", matchedMember)
		fmt.Println("\n\nlogged in member: ", memberID)
		if numberOfSharedBooks >= 2 {
			err = repository.DB.QueryRow("select username from til_member where id=?", matchedMember).Scan(&matchedUsername)
			if err != nil {
				fmt.Printf("Database lookup of username in HTTP request to %s failed with error: ", r.URL.Path, err)
			}

			fmt.Println("\n\nMatched with member: ", matchedUsername)
			_, err = repository.DB.Exec("insert into chatroom (book_title) ?")
			if err != nil {
				fmt.Printf("Database insert into chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
			}
		}
	}

	if len(books) == 0 {
		books = []Book{}
	}

	jsonBooks, _ := json.Marshal(&books)
	w.Write([]byte(jsonBooks))
}
