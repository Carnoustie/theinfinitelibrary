package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/argon2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

type Book struct {
	Id         int64  `json:id`
	MemberID   int64  `json:"member_id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	ChatRoomID int64  `json:"chatroom_id"`
}

type ChatPayLoad struct {
	Message    string `json:message`
	ChatRoomID string `json:chatroomid`
	Username   string `json:username`
}

var clients = make(map[http.ResponseWriter]chan string)
var chatRoomChannels = make(map[string][]chan string)
var mainChannel = make(chan ChatPayLoad)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //allow sharing response with client

	var bodyContents []byte
	bodyContents, err := io.ReadAll(r.Body)
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
	var salt []byte
	var pwHash []byte
	err = db.QueryRow("select salt,password_hash from til_member where username=?", u.Username).Scan(&salt, &pwHash)
	if err != nil {
		fmt.Printf("DB lookup of user in login failed with error: ", err)
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
	_, err = db.Exec("insert into til_member(username, salt, password_hash) values (?,?,?)", u.Username, salt, encryptedPassword)
	if err != nil {
		fmt.Println("Insert failed: ", err)
		w.Write([]byte("User already exists, pick a different username."))
	} else {
		w.Write([]byte("Welcome to The Infinite Library! :)"))
	}
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
	err = db.QueryRow("select id from til_member where username=?", u.Username).Scan(&userID)
	if err != nil {
		fmt.Printf("\n\nDatabase lookup of member %s failed with error %s\n\n", u.Username, err)
		return
	}
	_, err = db.Exec("insert into books (member_id, title, author) values(?,?,?)", userID, bk.Title, bk.Author)
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
	err := db.QueryRow("select id from til_member where username=?", u.Username).Scan(&memberID)
	if err != nil {
		fmt.Printf("Database lookup of member in HTTP request to %s failed with error: ", r.URL.Path, err)
		return
	}

	rows, err := db.Query("select id, title, author from books where member_id=?", memberID)
	if err != nil {
		fmt.Printf("Database lookup of book in HTTP request to %s failed with error: ", r.URL.Path, err)
		return
	}

	var books []Book
	var allBooks []Book
	var tempBook Book
	allBookRows, err := db.Query("select id, member_id, title, author from books")
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
				firstLevenshtein := levenshteinDistance(b_front, localbook_front, len(b_front), len(localbook_front))
				if firstLevenshtein < 2 {
					levenshteinDist = levenshteinDistance(strings.ReplaceAll(b.Title, " ", "_"), strings.ReplaceAll(localbook.Title, " ", "_"), len(b.Title), len(localbook.Title))
					fmt.Println("\n\nComputed Levenshtein: ", levenshteinDist)
					if levenshteinDist <= 2 {
						err = db.QueryRow("select username from til_member where id=?", localbook.MemberID).Scan(&localUserName)
						if err != nil {
							fmt.Printf("Database lookup of user %d in HTTP request to %s failed with error: ", localbook.MemberID, r.URL.Path, err)
						}

						countMatchesMap[localbook.MemberID]++
						if countMatchesMap[localbook.MemberID] >= 2 {
							err = db.QueryRow("select chatroom_id from books where id=?", b.Id).Scan(&checkChatroomID)
							if err != nil {
								fmt.Printf("Database lookup of chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
							}

							// fmt.Println("\n\nfetched: ", checkChatroomID)
							if checkChatroomID == 0 {
								_, err = db.Exec("insert into chatroom (book_title) values (?)", b.Title)

								if err != nil {
									fmt.Printf("Database insertion of book into chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
								} else {
									err = db.QueryRow("select id from chatroom where book_title = ?", b.Title).Scan(&chatRoomId)
									if err != nil {
										fmt.Printf("Database lookup of chatroom in HTTP request to %s failed with error: ", r.URL.Path, err)
									}

									_, err = db.Exec("update books set chatroom_id = ? where id = ? or id = ?", chatRoomId, b.Id, localbook.Id)
									if err != nil {
										fmt.Printf("Database update of books in HTTP request to %s failed with error: ", r.URL.Path, err)
									}
								}
							} else {
								_, err = db.Exec("update books set chatroom_id = ? where id = ?", chatRoomId, localbook.Id)
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

	fmt.Printf("\n\nCompleted chatroom assignment updates for books in Database in HTTP request to %s", r.URL.Path)
	rows, err = db.Query("select id, title, author, chatroom_id from books where member_id=?", memberID)
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
			err = db.QueryRow("select username from til_member where id=?", matchedMember).Scan(&matchedUsername)
			if err != nil {
				fmt.Printf("Database lookup of username in HTTP request to %s failed with error: ", r.URL.Path, err)
			}

			fmt.Println("\n\nMatched with member: ", matchedUsername)
			_, err = db.Exec("insert into chatroom (book_title) ?")
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

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow sharing response with client
	w.Header().Set("Content-Type", "application/json")

	chatId := strings.TrimPrefix(r.URL.Path, "/api/postMessage/")
	fmt.Println("\n\nEntered API endpoint %s", r.URL.Path)

	fmt.Println("\n\nhit postmessage api with chatid: ", chatId)
	var bodyContents []byte
	var payload ChatPayLoad
	bodyContents, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("\n\nParsing payload in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}

	err = json.Unmarshal(bodyContents, &payload)
	if err != nil {
		fmt.Printf("\n\nJSON parsing in HTTP request to %s failed with error %s\n\n", r.URL.Path, err)
		return
	}

	fmt.Println("\n\nRead from Browser: ", payload.Message, "\n\n")
	mainChannel <- payload

	fmt.Println("\n\nWrite complete\n\n")
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	chatId := strings.TrimPrefix(r.URL.Path, "/api/chatRoom/")

	fmt.Println("\n\nEntered chatroom api with chatid: ", chatId)

	clientChannel := make(chan string)
	chatRoomChannels[chatId] = append(chatRoomChannels[chatId], clientChannel)
	flush, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	var message string
	for {
		fmt.Println("\n\nEntered chat loop in chatroom api\n\n")
		message = <-clientChannel
		message = "data: " + message + "\n\n"
		fmt.Println("\n\nMessage to send:\n\n", message)
		fmt.Fprint(w, message)
		flush.Flush()
	}
}
