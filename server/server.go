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

	"github.com/google/uuid"
	"github.com/reply/types"
)

type User struct {
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Conn   net.Conn `json:"conn"`
}

var (
	mu             sync.Mutex
	ConnectedUsers = make(map[string]User)
)

func ValidateAction(jsonData []byte) (types.Header, error) {

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Invalid data received")
		return nil, fmt.Errorf("invalid data received")
	}

	action, ok := data["action"].(string)
	if !ok {
		return nil, fmt.Errorf("no action received")
	}

	switch action {
	case "TEXT_MESSAGE":
		var message types.Message
		err := json.Unmarshal(jsonData, &message)
		if err != nil {
			return nil, fmt.Errorf("invalid message data received")
		}
		return message, nil
	case "ABSENT":
		var absent types.Absent
		err := json.Unmarshal(jsonData, &absent)
		if err != nil {
			return nil, fmt.Errorf("invalid absent data received")
		}
		return absent, nil
	case "USER_JOIN":
		var status_update types.StatusUpdate
		err := json.Unmarshal(jsonData, &status_update)
		if err != nil {
			return nil, fmt.Errorf("invalid absent data received")
		}
		return status_update, nil
	default:
		return nil, fmt.Errorf("invalid default data received")
	}
}

func startListenting() net.Listener {
	listener, err := net.Listen("tcp", "localhost:6980")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	fmt.Println("Server started. Listening on localhost:6980")
	return listener
}

func receieveBegin(conn net.Conn, connectionId string) {
	beginJSON, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading client hash:", err)
	}
	message, err := ValidateAction([]byte(beginJSON))
	if err != nil {
		fmt.Println("Error validating message:", err)
	}
	mu.Lock()
	defer mu.Unlock()
	switch message.Type() {
	case "USER_JOIN":
		message := message.(types.StatusUpdate)
		fmt.Println("User: ", message.Name, " is now Online")
		cuser := User{
			Name:   message.Name,
			Status: message.Status,
			Conn:   conn,
		}
		ConnectedUsers[connectionId] = cuser
	default:
		log.Println("Client is sending Invalid Begin Struct")
	}
}

func createUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func ServerMain() {

	listener := startListenting()

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		connectionID := createUUID()

		receieveBegin(conn, connectionID)

		go handleClient(conn, connectionID)

	}
}

func handleClient(conn net.Conn, connectionID string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
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

		message, err := ValidateAction([]byte(jsonData))
		if err != nil {
			log.Fatalln("Error validating message:", err)
			continue
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
