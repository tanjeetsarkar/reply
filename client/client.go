package client

import (
	"bufio"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/reply/types"
)

func generateHash(n string) string {
	// generate sha512 hash
	hash := sha512.New()
	hash.Write([]byte(n))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func sanitzieUsername(n string) string {
	// remove spaces from username
	return strings.Replace(n, " ", "", -1)
}

func listenForMessages(clientHash string, conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		messageJSON, err := reader.ReadString('\n')
		// read the message from the server
		if err != nil {
			fmt.Println("Error reading message from server:", err)
			os.Exit(1)
		}

		// Validate the message
		message, err := ValidateAction([]byte(messageJSON))
		if err != nil {
			fmt.Println("Error validating message:", err)
			continue
		}

		switch message.Type() {
		case "TEXT_MESSAGE":
			message := message.(types.Message)
			fmt.Println(message.From[:5], ":", message.Message)
		case "ABSENT":
			message := message.(types.Absent)
			fmt.Println(message.SenderID[:5], "is absent")
		default:
			fmt.Println("Invalid message type")
		}
	}
}

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
	default:
		return nil, fmt.Errorf("invalid default data received")
	}
}

func ReplytoMessages(conn net.Conn, scanner *bufio.Scanner, clientHash string, recipientHash string, done chan<- bool) {
	for {
		fmt.Print(clientHash[:5], " : ")
		scanner.Scan()
		messageText := scanner.Text()

		if messageText == "/quit" {
			done <- true
			return
		}
		// Create and send the JSON message to the server
		message := types.Message{Action: "TEXT_MESSAGE", From: clientHash, To: recipientHash, Message: messageText}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshalling JSON message:", err)
			os.Exit(1)
		}

		_, err = conn.Write(append(messageJSON, '\n'))
		if err != nil {
			fmt.Println("Error sending JSON message to server:", err)
			os.Exit(1)
		}

	}
}

func clientInit(conn net.Conn) (net.Conn, bufio.Scanner, string, string) {

	fmt.Print("Enter your Name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	// Generate the client hash
	clientHash := generateHash(sanitzieUsername(name))
	fmt.Println("Your client hash:", clientHash)

	fmt.Println("Connected to server: ", conn.RemoteAddr())

	// Send the client hash to the server
	_, err := conn.Write([]byte(clientHash + "\n"))
	if err != nil {
		fmt.Println("Error sending client hash to server:", err)
		os.Exit(1)
	}

	fmt.Print("Enter the recipient's username: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	recipient := scanner.Text()
	recipientHash := generateHash(sanitzieUsername(recipient))

	return conn, *scanner, clientHash, recipientHash
}

func dialUp() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:6980")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ClientMain() {

	conn, err := dialUp()
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	defer conn.Close()

	done := make(chan bool)

	conn, scanner, clientHash, recipientHash := clientInit(conn)

	// Start a goroutine to listen for incoming messages
	go listenForMessages(clientHash, conn)

	// Start a goroutine to send messages
	go ReplytoMessages(conn, &scanner, clientHash, recipientHash, done)

	if <-done {
		return
	}
	// Block forever
	select {}

}
