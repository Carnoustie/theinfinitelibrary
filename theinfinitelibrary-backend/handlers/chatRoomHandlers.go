package handlers

import(
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ChatPayLoad struct {
	Message    string `json:message`
	ChatRoomID string `json:chatroomid`
	Username   string `json:username`
}

var err error;

var Clients = make(map[http.ResponseWriter]chan string)
var ChatRoomChannels = make(map[string][]chan string)
var MainChannel = make(chan ChatPayLoad)

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
	MainChannel <- payload

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
	ChatRoomChannels[chatId] = append(ChatRoomChannels[chatId], clientChannel)
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
