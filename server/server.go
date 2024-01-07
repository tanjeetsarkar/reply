package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/reply/types"
	validation "github.com/reply/util"
)

type User struct {
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Conn     net.Conn  `json:"conn"`
	LastSeen time.Time `json:"last_seen"`
}

func (u *User) updateTime() {
	u.LastSeen = time.Now()
	fmt.Println(u.Name, "last seen at", u.LastSeen)
}

var (
	mu             sync.Mutex
	ConnectedUsers = make(map[string]*User)
)

func startListenting() net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:6980")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	fmt.Println("Server started. Listening on localhost:6980")
	return listener
}

func receieveBegin(conn net.Conn, connectionId string) error {
	beginJSON, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading client hash:", err)
	}
	message, err := validation.ValidateAction([]byte(beginJSON))
	if err != nil {
		fmt.Println("Error validating message:", err)
	}
	mu.Lock()
	defer mu.Unlock()
	switch message.Type() {
	case "USER_JOIN":
		message := message.(types.StatusUpdate)
		fmt.Println("User: ", message.Name, " is now Online")
		fmt.Println("Joined at: ", time.Now())
		cuser := User{
			Name:     message.Name,
			Status:   message.Status,
			Conn:     conn,
			LastSeen: time.Now(),
		}
		ConnectedUsers[connectionId] = &cuser
		return nil
	default:
		log.Println("Client is sending Invalid Begin Struct")
		return fmt.Errorf("invalid begin struct")
	}
}

func createUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func ServerMain() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	listener := startListenting()

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		connectionID := createUUID()

		err = receieveBegin(conn, connectionID)
		if err != nil {
			continue
		}
		go handleClient(conn, connectionID)
	}
}

func handleClient(conn net.Conn, connectionID string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		go func() {
			ConnectedUsers[connectionID].updateTime()
		}()
		jsonData, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				mu.Lock()
				defer mu.Unlock()
				fmt.Println(ConnectedUsers[connectionID].Name, " got disconnected")
				delete(ConnectedUsers, connectionID)
			} else {
				fmt.Println("Error reading message from server:", err)
			}
			break
		}

		message, err := validation.ValidateAction([]byte(jsonData))
		if err != nil {
			log.Fatalln("Error validating message:", err)
			break
		}

		switch message.Type() {
		case "TEXT_MESSAGE":
			message := message.(types.Message)
			fmt.Println(message.From, ":", message.Message)
			go checkOnlineSend(conn, message)
		default:
			fmt.Println("Inavlid Message Recieved", message)
		}
	}
}

func marshalTextMessage(message types.Message) []byte {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON message:", err)
	}
	return messageJSON
}

func checkOnlineSend(conn net.Conn, message types.Message) {
	mu.Lock()
	defer mu.Unlock()
	for _, cuser := range ConnectedUsers {
		fmt.Println(cuser.Name)
		if message.To == cuser.Name {
			fmt.Println("sending", message.Message, "to", cuser.Name)
			_, err := cuser.Conn.Write(append([]byte(marshalTextMessage(message)), '\n'))
			if err != nil {
				log.Println("Error while sending to connected client", err)
				return
			}
			return
		}
	}
	err := sendUnavailable(conn, message.To)
	if err != nil {
		fmt.Println("cant send absent: ", err)
	}
}

func sendUnavailable(conn net.Conn, receiver string) error {

	absent := types.Absent{
		Action:   "ABSENT",
		SenderID: receiver,
	}

	absentJson, err := json.Marshal(absent)
	if err != nil {
		return err
	}

	_, err = conn.Write(append(absentJson, '\n'))
	return err
}
