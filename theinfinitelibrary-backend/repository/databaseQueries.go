package repository

import (
	"database/sql"
	"fmt"
)



type Book struct {
	Id         int64  `json:id`
	MemberID   int64  `json:"member_id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	ChatRoomID int64  `json:"chatroom_id"`
}

func GetUseridByUsername(username string) (int64, error) {
	var userID int64
	err := DB.QueryRow("select id from til_member where username=?", username).Scan(&userID)
	if err != nil {
		fmt.Printf("\n\nDatabase lookup of member %s failed with error %s\n\n", username, err)
		return -1, fmt.Errorf("\n\nDatabase lookup of member %s failed with error %s\n\n", username, err)
	}
	return userID, nil
}

func GetUsernameByID(userid int64) (string, error){
	var username string
	err := DB.QueryRow("select username from til_member where id=?", userid).Scan(&username)
	if err != nil {
		fmt.Printf("\n\nDatabase lookup of member %s failed with error %s\n\n", username, err)
		return "", fmt.Errorf("\n\nDatabase lookup of userid %d failed with error %s\n\n", userid, err)
	}
	return username, nil
}

func GetAllBooks() (*sql.Rows, error) {
	allBooks, err := DB.Query("select id, member_id, title, author, chatroom_id from books")
	if err != nil {
		fmt.Printf("\n\nDatabase fetch of all books failed with error %s\n\n", err)
		return nil, fmt.Errorf("\n\nDatabase fetch of all books failed with error %s\n\n", err)
	}
	return allBooks, nil
}

func GetBooksByMemberID(memberID int64) (*sql.Rows, error) {
	Books, err := DB.Query("select id, member_id, title, author, chatroom_id from books where member_id=?", memberID)
	if err != nil {
		fmt.Printf("\n\nDatabase fetch of books with member id %d failed with error %s\n\n", memberID, err)
		return nil, fmt.Errorf("\"\n\nDatabase fetch of books with member id %d failed with error %s\n\n", memberID, err)
	}
	return Books, nil
}


func PersistBook(memberID int64, title string, author string) (error){
	_, err := DB.Exec("insert into books (member_id, title, author) values(?,?,?)", memberID, title, author)
	if err!=nil{
		return fmt.Errorf("\n\nDatabase insert of book %d, %s, %s failed with error %s", memberID, title, author, err)
	}
	return nil
}

func GetChatroomIdFromBookId(bookID int64) (int64, error){
	var chatRoomID int64
	err := DB.QueryRow("select chatroom_id from books where id=?", bookID).Scan(&chatRoomID)
	if err!=nil{
		return -1, fmt.Errorf("\n\ndatabase lookup of chatroom for book with id %d failed with error %s", bookID, err)
	}
	return chatRoomID, nil
}

func GetChatroomIdFromBookTitle(booktitle string) (int64, error){
	var chatRoomID int64
	err := DB.QueryRow("select id from chatroom where book_title = ?", booktitle).Scan(&chatRoomID)
	if err!=nil{
		return -1, fmt.Errorf("\n\ndatabase lookup of chatroom for book with title %s failed with error %s", booktitle, err)
	}
	return chatRoomID, nil
}

func AddBookToChatroom(booktitle string) (error){
	_, err:= DB.Exec("insert into chatroom (book_title) values (?)", booktitle)
	if err!=nil{
		return fmt.Errorf("\n\ndatabase insert of book with title %s into chatroom table failed with error%s", booktitle, err)
	}
	return nil
}

func SetChatroomIdFromMultipleBooks(bookId1 int64, bookId2 int64, chatRoomID int64) (error){
	_, err := DB.Exec("update books set chatroom_id = ? where id = ? or id = ?", chatRoomID, bookId1, bookId2)
	if err!=nil{
		return fmt.Errorf("\n\ndatabase update of chatroomid %d with for books %d and %d failed with error%s", chatRoomID, bookId1, bookId2, err)
	}
	return nil
}

func SetChatroomIdFromSingleBook(bookID int64, chatRoomID int64) (error){
	_, err := DB.Exec("update books set chatroom_id = ? where id = ?", chatRoomID, bookID)
	if err!=nil{
		return fmt.Errorf("\n\ndatabase update of chatroomid %d with for book with id %d failed with error%s", chatRoomID, bookID, err)
	}
	return nil
}
