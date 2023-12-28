package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/reply/client"
	"github.com/reply/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WritePush struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func Wp(conn *websocket.Conn, writePump chan types.Message, from string) {
	defer close(writePump) // Defer closing the writePump channel
	for {
		_, messageText, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error", err)
			break
		}
		var wp WritePush
		err = json.Unmarshal(messageText, &wp)
		if err != nil {
			log.Println("Unmarshal error", err)
			continue
		}
		clientHash := client.SanitzieUsername(from)
		recipientHash := client.SanitzieUsername(wp.To)
		message := types.Message{Action: "TEXT_MESSAGE", From: clientHash, To: recipientHash, Message: wp.Message}
		// implement message unmarshal here to create type.Message.
		writePump <- message
	}
}

func Rp(conn *websocket.Conn, readPump chan string) {
	defer close(readPump) // Defer closing the readPump channel
	for reads := range readPump {
		err := conn.WriteMessage(websocket.TextMessage, []byte(reads))
		if err != nil {
			log.Println("Writing error", err)
			break
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	var (
		writePump = make(chan types.Message)
		readPump  = make(chan string)
		done      = make(chan bool)
	)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	from := r.URL.Query().Get("from")
	defer conn.Close()
	go client.ClientMain(writePump, readPump, from, done)
	go Wp(conn, writePump, from)
	go Rp(conn, readPump)
	select {}
}
