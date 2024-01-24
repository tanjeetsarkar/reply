package clientv2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/reply/types"
	validation "github.com/reply/util"
)

type ClientV2 struct {
	WritePump  chan types.Message
	ReadPump   chan types.Message
	ActivePump chan types.StatusResponse
	net.Conn
	Hostname string
	ActiveTo string
	ClientID string
	Reader   *bufio.Reader
}

func NewClientV2(clientId string,
	writePump chan types.Message,
	readPump chan types.Message,
	ActivePump chan types.StatusResponse,
	tcpAddr string) *ClientV2 {
	conn, err := net.Dial("tcp", tcpAddr)

	if err != nil {
		log.Fatalln("Error Connecting to Server: ", err)
	}
	log.Println("Client connected to Server ... ")
	identifier, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

	c := &ClientV2{
		WritePump: writePump,
		ReadPump:  readPump,
		Conn:      conn,
		Hostname:  identifier,
		ClientID:  clientId,
		Reader:    bufio.NewReader(conn),
	}

	return c
}

func (c *ClientV2) SendAuth() {
	userJoin := types.StatusUpdate{
		Action: "USER_JOIN",
		Name:   c.ClientID,
		Status: "ONLINE",
	}
	userJoinJSON, err := json.Marshal(userJoin)

	if err != nil {
		log.Fatalln("Error marshalling JSON message:", err)
	}
	c.Write([]byte(append(userJoinJSON, '\n')))

}

func (c *ClientV2) CheckLastSeen(checkClientID string) {
	checkStatus := types.CheckStatus{
		Action: "CHECK_STATUS",
		Chash:  checkClientID,
	}

	checkStatusJSON, err := json.Marshal(checkStatus)

	if err != nil {
		log.Println("Error marshalling JSON message:", err)
	}

	_, err = c.Write([]byte(append(checkStatusJSON, '\n')))

	if err != nil {
		log.Println("Error writing Check Status to server : ", err)
	}
}

func (c *ClientV2) SendTextMessage() {
	fmt.Println("Sending Pending Messages ... ")
	for {
		for message := range c.WritePump {
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Fatalln("Error Marshalling message JSON", err)
			}
			_, err = c.Write([]byte(append(messageJSON, '\n')))
			if err != nil {
				log.Fatalln("Error sending Writing message to server", err)
			}
		}
	}
}

func (c *ClientV2) ListenForMessages() {
	for {
		messageJSON, err := c.Reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message from server:", err)
			if err.Error() == "EOF" {
				log.Fatalln("Server closed connection")
			}
		}
		message, err := validation.ValidateAction([]byte(messageJSON))
		if err != nil {
			log.Println("Error validating message:", err)
		}

		switch message.Type() {
		case "TEXT_MESSAGE":
			message := message.(types.Message)
			c.ReadPump <- message
		case "ABSENT":
			message := message.(types.Absent)
			fmt.Println(message.SenderID, "is offline")
		case "STATUS_RESPONSE":
			message := message.(types.StatusResponse)
			c.ActivePump <- message
			fmt.Println(message.Chash, "Last Seen at: ", message.LastSeen)
		default:
			fmt.Println("Invalid message type")
		}
	}
}

func (c *ClientV2) IOloop(controlChan <-chan struct{}) {
	select {
	case <-controlChan:
		log.Println("Stopping Client loop..")
		return
	default:
		go c.SendTextMessage()
		go c.ListenForMessages()
	}
}
