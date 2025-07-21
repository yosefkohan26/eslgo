package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command"
)

func main() {
	// Set up a channel to handle shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	conn, err := eslgo.Dial("146.190.77.119:8021", "DialB2025!", func() {
		log.Println("Connection disconnected.")
	})
	if err != nil {
		log.Fatalf("Error connecting to FreeSWITCH: %v", err)
	}
	defer conn.ExitAndClose()

	log.Println("Successfully connected to FreeSWITCH")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Enable JSON events
	_, err = conn.SendCommand(ctx, command.Event{Format: "json", Listen: []string{"all"}})
	if err != nil {
		log.Fatalf("Error enabling JSON events: %v", err)
	}
	log.Println("Enabled JSON event listener")

	// 2. Enable logs at the INFO level (level 6)
	_, err = conn.EnableLogs(ctx, 6)
	if err != nil {
		log.Fatalf("Error enabling logs: %v", err)
	}
	log.Println("Enabled log reception at INFO level")

	// 3. Start a goroutine to listen on the log channel
	go func() {
		logChannel := conn.LogChannel()
		for {
			select {
			case logEntry, ok := <-logChannel:
				if !ok {
					log.Println("Log channel closed.")
					return
				}
				fmt.Println("--- LOG ENTRY RECEIVED ---")
				fmt.Printf("Level: %d\n", logEntry.Level())
				fmt.Printf("Message: %s\n", logEntry.Message())
				fmt.Println("-------------------------")
			case <-ctx.Done():
				return
			}
		}
	}()

	// 4. Register an event listener for all events
	conn.RegisterEventListener(eslgo.EventListenAll, func(event *eslgo.Event) {
		fmt.Println("--- JSON EVENT RECEIVED ---")
		fmt.Printf("Event Name: %s\n", event.GetName())
		fmt.Printf("Unique-ID: %s\n", event.GetHeader("Unique-ID"))
		fmt.Printf("Raw Body: %s\n", string(event.Body))
		fmt.Println("--------------------------")
	})

	log.Println("Listening for events and logs. Press Ctrl+C to exit.")

	// Wait for a shutdown signal
	<-shutdown
	log.Println("Shutdown signal received, closing connection.")
}