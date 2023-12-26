package main

import (
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

func Wp(conn *websocket.Conn, writePump chan types.Message) {
	defer close(writePump) // Defer closing the writePump channel
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error", err)
			break
		}
		// implement message marshal here to create type.Message.
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
		from      string
		done      = make(chan bool)
	)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	from = r.URL.Query().Get("from")
	defer conn.Close()
	go client.ClientMain(writePump, readPump, from, done)
	go Wp(conn, writePump)
	go Rp(conn, readPump)
	select {}
}
