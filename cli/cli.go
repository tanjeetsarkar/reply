package main

import (
	"fmt"
	"log"
	"os"

	"github.com/reply/clientv2"
	"github.com/reply/server"
	"github.com/reply/types"
)

type UserAuth struct {
	Name     string
	Password string
}

func (u *UserAuth) isAuthentic() bool {
	//TODO: Implement authentication strategy
	return u.Name == u.Password
}

func authenticate() bool {
	// fmt.Print("username: ")
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// name := scanner.Text()
	name := "Tanjeet"
	user := UserAuth{
		Name:     name,
		Password: name,
	}
	return user.isAuthentic()
}

func startClientV2() {
	writePump := make(chan types.Message)
	readPump := make(chan types.Message)
	tcpAddr := "192.168.0.105:6980"
	clientId := os.Args[3]

	c := clientv2.NewClientV2(clientId, writePump, readPump, tcpAddr)
	c.SendAuth()
	go func() {

		fmt.Println("pushing sample to writepump")
		writePump <- types.Message{
			Action:  "TEXT_MESSAGE",
			From:    clientId,
			To:      "sum",
			Message: "Test Message",
		}
	}()
	c.SendTextMessage()
	fmt.Println("Sent sample to writepump")
	fmt.Println(c.Hostname)
	fmt.Println(c.ClientID)
	close(writePump)
	defer c.Close()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage ... reply client start")
	}

	if os.Args[1] == "server" {
		if os.Args[2] == "start" {
			log.Println("starting Server ... ")
			if authenticate() {
				server.ServerMain()
			} else {
				log.Fatalln("Authentication Failed")
			}
		}
	}
	if os.Args[1] == "client" {
		if os.Args[2] == "start" {
			log.Println("starting Client ... ")
			if authenticate() {
				startClientV2()
			} else {
				log.Fatalln("Authentication Failed")
			}
		}
	}
}
