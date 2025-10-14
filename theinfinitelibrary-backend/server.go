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
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/argon2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

type Book struct {
	Id       int64  `json:id`
	MemberID int64  `json:"member_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
}

var clients = make(map[http.ResponseWriter]chan string)
var mainChannel = make(chan string)

func main() {
	fmt.Println("\n\n\n\nThe Infinite Library Server is running\n\n\n\n")

	db_username := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")

	go broadcaster(mainChannel, clients)
	// go spawnChatRoom()

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

		fmt.Println("\n\nUsername: ", u.Username)
		if err != nil {
			fmt.Println("DB lookup failed with error: ", err)
			_, _ = w.Write([]byte("error"))
			return
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

	http.HandleFunc("/api/addbook", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client

		var bodyContents []byte
		bodyContents, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("\n\nbody: ", string(bodyContents))

		var bk Book
		var u User
		err = json.Unmarshal(bodyContents, &bk)
		err = json.Unmarshal(bodyContents, &u)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n\nTitle: ", bk.Title, "\n\nAuthor: ", bk.Author)

		var userID int
		err = db.QueryRow("select id from til_member where username=?", u.Username).Scan(&userID)
		if err != nil {
			fmt.Println("db lookup 1 failed with error: ", err)
		}

		fmt.Println("\n\nQuery result: ", int64(userID))

		_, err = db.Exec("insert into books (member_id, title, author) values(?,?,?)", userID, bk.Title, bk.Author)
		if err != nil {
			fmt.Println("db insert 2 failed with", err)
		}
	})

	http.HandleFunc("/api/getbooks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client
		w.Header().Set("Content-Type", "application/json")

		var bodyContents []byte
		bodyContents, err = io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var u User

		err = json.Unmarshal(bodyContents, &u)
		if err != nil {
			panic(err)
		}
		var memberID int64
		err := db.QueryRow("select id from til_member where username=?", u.Username).Scan(&memberID)
		if err != nil {
			fmt.Println("db lookup 3 failed with error: ", err)
		}
		// var result []string

		// err = db.QueryRow("select * from books where member_id=?", memberID).Scan(&result)
		// row := db.QueryRow("select * from books where member_id=?", memberID)

		rows, _ := db.Query("select id, title, author from books where member_id=?", memberID)

		var books []Book

		var allBooks []Book

		var tempBook Book
		allBookRows, _ := db.Query("select id, member_id, title, author from books")
		for allBookRows.Next() {
			allBookRows.Scan(&tempBook.Id, &tempBook.MemberID, &tempBook.Title, &tempBook.Author)
			allBooks = append(allBooks, tempBook)
			//fmt.Println("\n\nJust added: ", tempBook.Title)
		}

		var b Book

		//var matchedUsersIDs []int

		var levenshteinDist int

		var localUserName string

		countMatchesMap := make(map[int64]int)

		var chatRoomId int64

		var checkChatroomID int64

		for rows.Next() {
			rows.Scan(&b.Id, &b.Title, &b.Author)

			for _, localbook := range allBooks {

				if localbook.MemberID != memberID {
					// fmt.Println("\n\nlocalBook member ID: ", localbook.MemberID)
					// fmt.Println("\n\nlogged in member ID: ", memberID)
					b_front := strings.Split(b.Title, " ")[0]
					localbook_front := strings.Split(localbook.Title, " ")[0]
					firstLevenshtein := levenshteinDistance(b_front, localbook_front, len(b_front), len(localbook_front))
					if firstLevenshtein < 2 {
						levenshteinDist = levenshteinDistance(strings.ReplaceAll(b.Title, " ", "_"), strings.ReplaceAll(localbook.Title, " ", "_"), len(b.Title), len(localbook.Title))
						fmt.Println("\n\nComputed Levenshtein: ", levenshteinDist)
						if levenshteinDist <= 2 {
							err = db.QueryRow("select username from til_member where id=?", localbook.MemberID).Scan(&localUserName)
							if err != nil {
								fmt.Println("\n\nDB lookup 5 failed with error: ", err)
							}
							// fmt.Println("\n\nMatched member: ", localUserName)
							// fmt.Println("\n\nb.Title: ", b.Title)
							// fmt.Println("\n\nlocalBook.Title: ", localbook.Title)
							countMatchesMap[localbook.MemberID]++
							if countMatchesMap[localbook.MemberID] >= 2 {
								err = db.QueryRow("select chatroom_id from books where id=?", b.Id).Scan(&checkChatroomID)
								if err != nil {
									fmt.Println("\n\nfetching chatroom id failed with error: ", err)
								}
								// fmt.Println("\n\nfetched: ", checkChatroomID)
								if checkChatroomID == 0 {
									_, err = db.Exec("insert into chatroom (book_title) values (?)", b.Title)
									if err != nil {
										fmt.Println("\n\ninserting into chatroom failed with error: ", err)
									} else {
										err = db.QueryRow("select id from chatroom where book_title = ?", b.Title).Scan(&chatRoomId)
										if err != nil {
											fmt.Println("\n\nfetch chatroom id failed with error: ", err)
										}
										_, err = db.Exec("update books set chatroom_id = ? where id = ? or id = ?", chatRoomId, b.Id, localbook.Id)
										if err != nil {
											fmt.Println("\n\nupdating chatroom id's into books failed with error: ", err)
										}
									}
								} else {
									_, err = db.Exec("update books set chatroom_id = ? where id = ?", chatRoomId, localbook.Id)
									if err != nil {
										fmt.Println("\n\nupdating chatroom id's into books failed with error: ", err)
									}
								}
							}
						}
					}
				}
			}

			books = append(books, b)
			fmt.Println("\n\nbooktitle: ", b.Title)

			//fmt.println("\n\nlocalbook: ", b.author, "   ", b.title)
		}

		var matchedUsername string
		for matchedMember, numberOfSharedBooks := range countMatchesMap {
			fmt.Println("\n\nsharedBookCount: ", numberOfSharedBooks)
			fmt.Println("\n\nwith member: ", matchedMember)
			fmt.Println("\n\nlogged in member: ", memberID)
			if numberOfSharedBooks >= 2 {
				_ = db.QueryRow("select username from til_member where id=?", matchedMember).Scan(&matchedUsername)
				fmt.Println("\n\nMatched with member: ", matchedUsername)
				_, err = db.Exec("insert into chatroom (book_title) ?")
			}
		}

		//var tempStoreUser User

		//// var username string

		//var matchedBooks []Book

		//for rows.next() {
		//	rows.scan(&b.title, &b.author)

		//	books = append(books, b)
		//	fmt.println("\n\nbooktitle: ", b.title)

		//	matchrows, err := db.query("select username from til_member where id in (select member_id from books where title=?)", b.title)

		//	if err == nil {
		//		for matchrows.next() {
		//			matchrows.scan(&tempstoreuser.username)
		//			if tempstoreuser.username != u.username {
		//				matchedusers = append(matchedusers, tempstoreuser)
		//				fmt.println("\n\nmatched user: ", tempstoreuser.username)
		//			}
		//		}
		//	}

		//	//fmt.println("\n\nlocalbook: ", b.author, "   ", b.title)
		//}

		if len(books) == 0 {
			books = []Book{}
		}

		jsonBooks, _ := json.Marshal(&books)
		// jsonMatchedUsers, _ := json.Marshal(&matchedUsers)

		// fmt.Println("\n\njsonbooks: ", string(jsonBooks))
		// fmt.Println("\n\njsonbooks: ", string(jsonMatchedUsers))

		w.Write([]byte(jsonBooks))
		// w.Write([]byte(jsonMatchedUsers))

		// if err != nil {
		// 	fmt.Println("db lookup 4 failed with error: ", err)
		// }

		fmt.Println("\n\nmemberID: ", memberID)

		// fmt.Println("\n\nBookarray: ", result)

		// fmt.Println("\n\nrow: ", row)
	})

	http.HandleFunc("/api/postMessage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow sharing response with client
		w.Header().Set("Content-Type", "application/json")

		fmt.Println("\n\nHit here")

		var bodyContents []byte
		bodyContents, err = io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n\nRead from Browser: ", string(bodyContents), "\n\n")
		mainChannel <- string(bodyContents)
		fmt.Println("\n\nWrite complete\n\n")

	})

	http.HandleFunc("/api/chatRoom", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //Allow sharing response with client
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		clientChannel := make(chan string)

		clients[w] = clientChannel

		//clients = append(clients, w)

		//var bodyContents []byte
		flush, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, ": init\n\n")
		flush.Flush()
		var message string
		for {
			fmt.Println("\n\nHit chat loop\n\n")
			message = <-clientChannel
			message = "data: " + message + "\n\n"
			//message += "\n\n"
			fmt.Println("\n\nMessage to send:\n\n", message)
			//time.Sleep(1 * time.Second)
			// io.ReadAll(r, []byte())
			fmt.Fprint(w, message)
			//fmt.Fprint(w, "data: Bonjour\n\n")
			flush.Flush()
			//fmt.Println("data: ticked\n\n")
		}
	})

	http.ListenAndServe(":8080", nil)
}

func broadcaster(ch chan string, clients map[http.ResponseWriter]chan string) {
	var chatMessage string

	for {
		chatMessage = <-ch
		fmt.Println("\n\nRead following message: ", chatMessage, "\n\n")
		for _, channel := range clients {
			fmt.Println("\n\nhit in broadcast loop\n\n")
			channel <- chatMessage
			//fmt.Fprint(client, chatMessage)
			//client.Write([]byte(chatMessage))
		}
	}
}

func levenshteinDistance(s1 string, s2 string, len_s1 int, len_s2 int) int {
	// fmt.Println("\n\ns1: ", s1, "\n\ns2: ", s2, "\n\n")
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
