package main

import (
	"context"
	"fmt"

	"github.com/reply/client"
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
	go func() {
		done <- true
	}()
	close(done)
	close(writePump)
	close(readPump)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

var (
	writePump = make(chan string)
	readPump  = make(chan string)
	done      = make(chan bool)
)

func (a *App) Start_client(from string, to string) {
	go func() {
		done <- false
	}()
	go client.ClientMain(writePump, readPump, from, to, done)
	runtime.EventsEmit(a.ctx, "clientStarted", "Client Started")
	go func() {
		for {
			for reads := range readPump {
				runtime.EventsEmit(a.ctx, "recieveMessage", reads)
			}
		}
	}()
}

func (a *App) SendMessage(message string) string {
	writePump <- message
	return "Message sent"
}
