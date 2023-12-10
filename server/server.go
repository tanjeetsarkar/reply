package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/reply/types"
)

// BeginStruct represents the structure of the "begin" message
type BeginStruct struct {
	Username string `json:"username"`
}

// Client struct represents a connected client
type Client struct {
	conn net.Conn
	hash string
}

// Channel struct represents a private communication channel

var clients = make(map[string]Client)

func handleClient(client Client) {
	defer client.conn.Close()

	reader := bufio.NewReader(client.conn)
	for {
		// Read the JSON message from the client
		jsonData, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading JSON message:", err)
			break
		}
		fmt.Println("Received JSON message: \n", jsonData)

		var message types.Message
		err = json.Unmarshal([]byte(jsonData), &message)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			break
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshalling JSON message:", err)
			continue
		}

		val, ok := clients[message.To]
		if !ok {
			err := sendUnavailable(client, message.To)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Reciever Unavailable")
			continue
		}
		// Send the message to the TO client
		_, err = val.conn.Write(append(messageJSON, '\n'))
		if err != nil {
			fmt.Println("Error sending message to :", err)
			err := sendUnavailable(client, message.To)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Reciever Unavailable")
		}
	}
}

func sendUnavailable(c Client, receiver string) error {

	absent := types.Absent{
		Action:   "ABSENT",
		SenderID: receiver,
	}

	absentJson, err := json.Marshal(absent)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(append(absentJson, '\n'))
	return err
}

func ServerMain() {
	listener, err := net.Listen("tcp", "localhost:6980")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server started. Listening on localhost:6980")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Read the client hash from the client
		hash, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading client hash:", err)
			continue
		}
		hash = hash[:len(hash)-1] // Remove the newline character
		fmt.Println("Client connected:", hash)

		client := Client{conn: conn, hash: hash}
		clients[hash] = client

		// go handleBegin(client)
		go handleClient(client)
	}
}
