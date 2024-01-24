package main

import (
	"context"
	"fmt"

	"github.com/reply/client"
	"github.com/reply/clientv2"
	"github.com/reply/models"
	"github.com/reply/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	fmt.Println("Shutting down...")
	close(done)
	close(writePump)
	close(readPump)
}

var (
	writePump  = make(chan types.Message)
	readPump   = make(chan types.Message)
	activePump = make(chan types.StatusResponse)
	done       = make(chan struct{})
)

const TCPADDR = "192.168.0.105:6980"

func (a *App) Start_client(clientID string) {
	c := clientv2.NewClientV2(clientID, writePump, readPump, activePump, TCPADDR)
	c.SendAuth()
	c.IOloop(done)
	runtime.EventsEmit(a.ctx, "clientStarted", "Client Started")
	go func() {
		for reads := range readPump {
			runtime.EventsEmit(a.ctx, "recieveMessage", reads.Message)
		}
	}()
}

func (a *App) SendMessage(messageText string, from string, to string) string {
	clientHash := client.SanitzieUsername(from)
	recipientHash := client.SanitzieUsername(to)
	message := types.Message{Action: "TEXT_MESSAGE", From: clientHash, To: recipientHash, Message: messageText}
	writePump <- message
	return "Message sent"
}

func (a *App) AddContact(name string, uid string) bool {

	contact := models.ContactsList{
		Name: name,
		Uid:  uid,
	}

	return models.ContactsInsert(contact)

}

func (a *App) GetContacts() []models.ContactsList {
	return models.ContactsGetAll()
}
