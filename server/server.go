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

	// _ "net/http/pprof"

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
}

func (u *User) CheckuserPresence(chash string) bool {
	return u.Name == chash
}

func (u *User) CheckLastSeen() time.Time {
	return u.LastSeen
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

func receieveBegin(conn net.Conn) (chash string, err error) {
	beginJSON, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading client hash:", err)
	}
	message, err := validation.ValidateAction([]byte(beginJSON))
	if err != nil {
		fmt.Println("Error validating message:", err)
		return "", err
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
		ConnectedUsers[message.Name] = &cuser
		return message.Name, nil
	default:
		log.Println("Client is sending Invalid Begin Struct")
		return "", fmt.Errorf("invalid begin struct")
	}
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

		chash, err := receieveBegin(conn)
		if err != nil {
			conn.Close()
			continue
		}
		go handleClient(conn, chash)
	}
}

func handleClient(conn net.Conn, chash string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		ConnectedUsers[chash].updateTime()
		jsonData, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				mu.Lock()
				fmt.Println(ConnectedUsers[chash].Name, " got disconnected")
				mu.Unlock()
				delete(ConnectedUsers, chash)
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
		case "CHECK_STATUS":
			message := message.(types.CheckStatus)
			fmt.Println("Checking status of : ", message.Chash)
			mu.Lock()
			if user, ok := ConnectedUsers[message.Chash]; ok {
				fmt.Println("User", user.Name, "is online")
				statusResponse := types.StatusResponse{
					Action:   "STATUS_RESPONSE",
					Chash:    message.Chash,
					LastSeen: user.LastSeen,
				}
				_, err := conn.Write(append(marshalStatusResponse(statusResponse), '\n'))
				if err != nil {
					fmt.Println("Error sending status response:", err)
				}
			} else {
				fmt.Println("User is offline")
				err := sendUnavailable(conn, message.Chash)
				if err != nil {
					fmt.Println("cant send absent: ", err)
				}
			}
			mu.Unlock()

		default:
			fmt.Println("Inavlid Message Recieved", message)
		}
	}
}

func marshalStatusResponse(message types.StatusResponse) []byte {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling JSON message:", err)
	}
	return messageJSON
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
	// for _, cuser := range ConnectedUsers {
	// 	fmt.Println(cuser.Name)
	// 	if message.To == cuser.Name {
	// 		fmt.Println("sending", message.Message, "to", cuser.Name)
	// 		_, err := cuser.Conn.Write(append([]byte(marshalTextMessage(message)), '\n'))
	// 		if err != nil {
	// 			log.Println("Error while sending to connected client", err)
	// 			return
	// 		}
	// 		return
	// 	}
	// }
	if user, ok := ConnectedUsers[message.To]; ok {
		fmt.Println("sending", message.Message, "to", user.Name)
		_, err := user.Conn.Write(append([]byte(marshalTextMessage(message)), '\n'))
		if err != nil {
			log.Println("Error while sending to connected client", err)
			return
		}
		return
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
