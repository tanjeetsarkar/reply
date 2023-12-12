package main

import (
	"fmt"
	"log"
	"os"

	"github.com/reply/server"
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage ... reply client start")
		return
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
}
