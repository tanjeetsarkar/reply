package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/reply/types"
	validation "github.com/reply/util"
)

func SanitzieUsername(n string) string {
	// remove spaces from username
	return strings.Replace(n, " ", "", -1)
}

func listenForMessages(clientHash string, conn net.Conn, absentQ chan string, readPump chan string, done chan<- bool, nac *ActiveChat) {
	reader := bufio.NewReader(conn)

	for {
		CheckStatus(conn, nac)
		messageJSON, err := reader.ReadString('\n')
		// read the message from the server
		if err != nil {
			fmt.Println("Error reading message from server:", err)
			if err.Error() == "EOF" {
				fmt.Println("Server closed connection")
				go func() {
					done <- true
				}()
				return
			}
		}

		// Validate the message
		message, err := validation.ValidateAction([]byte(messageJSON))
		if err != nil {
			fmt.Println("Error validating message:", err)
			continue
		}

		switch message.Type() {
		case "TEXT_MESSAGE":
			message := message.(types.Message)
			fmt.Println(message.From, ":", message.Message)
			readPump <- message.Message
		case "ABSENT":
			message := message.(types.Absent)
			fmt.Println(message.SenderID, "is absent", message)
			go func() {

				absentQ <- message.SenderID
			}()
		case "STATUS_RESPONSE":
			message := message.(types.StatusResponse)
			fmt.Println(message.Chash, "Last Seen at: ", message.LastSeen)
		default:
			fmt.Println("Invalid message type")
		}
	}
}

func ReplytoMessages(
	conn net.Conn,
	clientHash string,
	done chan<- bool,
	msgQ chan MessageQueue,
	absentQ chan string,
	writePump chan types.Message,
	nac *ActiveChat,
) {
	for {
		for message := range writePump {
			nac.SetRhash(message.To)

			go sendPendingMessages(msgQ, conn, absentQ)
			fmt.Print(clientHash, " : ")
			// scanner.Scan()

			fmt.Println("SENDING", message.Message)

			// messageText := scanner.Text()

			if message.Message == "/quit" {
				go func() {
					done <- true
				}()
				return
			}
			// Create and send the JSON message to the server
			// message := types.Message{Action: "TEXT_MESSAGE", From: clientHash, To: recipientHash, Message: messageText}

			messageJSON, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error marshalling JSON message:", err)
				os.Exit(1)
			}

			msgQ <- MessageQueue{
				msgP: messageJSON,
			}
		}
	}
}

func sendPendingMessages(msgQ chan MessageQueue, conn net.Conn, absentQ chan string) {
	for msg := range msgQ {
		go SendToServer(conn, msg.msgP)
	}
}

func SendToServer(conn net.Conn, messageJSON []byte) {

	_, err := conn.Write(append(messageJSON, '\n'))
	if err != nil {
		fmt.Println("Error sending JSON message to server:", err)
		os.Exit(1)
	}
}

func clientInit(conn net.Conn, from string) (net.Conn, string) {

	// fmt.Print("Enter your Name: ")
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	name := from

	// Generate the client hash
	clientHash := SanitzieUsername(name)
	fmt.Println("Your client hash:", clientHash)

	fmt.Println("Connected to server: ", conn.RemoteAddr())

	userJoin := types.StatusUpdate{
		Action: "USER_JOIN",
		Name:   clientHash,
		Status: "ONLINE",
	}

	userJoinJSON, err := json.Marshal(userJoin)

	if err != nil {
		fmt.Println("Error marshalling JSON message:", err)
		os.Exit(1)
	}

	// Send the client hash to the server
	_, err = conn.Write([]byte(append(userJoinJSON, '\n')))
	if err != nil {
		fmt.Println("Error sending client hash to server:", err)
		os.Exit(1)
	}

	// scanner = bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// recipient := scanner.Text()

	return conn, clientHash
}

func dialUp() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:6980")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ClientMain(writePump chan types.Message, readPump chan string, from string, done chan bool) {

	conn, err := dialUp()
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	defer conn.Close()

	msgQ := make(chan MessageQueue)
	absentQ := make(chan string)

	conn, clientHash := clientInit(conn, from)

	nac := NewActiveChat()

	// Start a goroutine to listen for incoming messages
	go listenForMessages(clientHash, conn, absentQ, readPump, done, nac)

	// Start a goroutine to send messages
	go ReplytoMessages(conn, clientHash, done, msgQ, absentQ, writePump, nac)

	if <-done {
		close(msgQ)
		close(absentQ)
		return
	}
	// Block forever
	select {}
}

type MessageQueue struct {
	msgP []byte
}

// func NewMessageQueue() *MessageQueue {
// 	return &MessageQueue{}
// }

type ActiveChat struct {
	Rhash string
}

func NewActiveChat() *ActiveChat {
	return &ActiveChat{}
}

func (ac *ActiveChat) SetRhash(rhash string) {
	ac.Rhash = rhash
}

func CheckStatus(conn net.Conn, ac *ActiveChat) {

	fmt.Println("Checking status of", ac.Rhash)

	if ac.Rhash != "" {

		checkStatus := types.CheckStatus{
			Action: "CHECK_STATUS",
			Chash:  ac.Rhash,
		}

		checkStatusJSON, err := json.Marshal(checkStatus)

		if err != nil {
			fmt.Println("Error marshalling JSON message:", err)
			os.Exit(1)
		}

		SendToServer(conn, checkStatusJSON)
	}

}
