package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/reply/client"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Wp(conn *websocket.Conn, writePump chan string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("REading error", err)
			break
		}
		writePump <- string(message)
	}
}

func Rp(conn *websocket.Conn, readPump chan string) {
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
		writePump chan string
		readPump  chan string
		from      string
		to        string
	)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	go client.ClientMain(writePump, readPump, from, to)

	go Wp(conn, writePump)
	go Rp(conn, readPump)
}
