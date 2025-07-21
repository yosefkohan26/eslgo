package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command"
)

func main() {
	// Set up a channel to handle shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	conn, err := eslgo.Dial("146.190.77.119:8021", "DialB2025!", func() {
		log.Println("✓ Connection disconnected cleanly")
	})
	if err != nil {
		log.Fatalf("Error connecting to FreeSWITCH: %v", err)
	}
	defer conn.ExitAndClose()

	log.Println("✓ Successfully connected to FreeSWITCH")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Enable JSON events
	_, err = conn.SendCommand(ctx, command.Event{Format: "json", Listen: []string{"HEARTBEAT", "RE_SCHEDULE"}})
	if err != nil {
		log.Fatalf("Error enabling JSON events: %v", err)
	}
	log.Println("✓ Enabled JSON event listener")

	// 2. Enable logs at WARNING level (level 4)
	_, err = conn.EnableLogs(ctx, 4)
	if err != nil {
		log.Fatalf("Error enabling logs: %v", err)
	}
	log.Println("✓ Enabled log reception at WARNING level")

	// Track received items
	eventCount := 0
	logCount := 0

	// 3. Start a goroutine to listen on the log channel
	go func() {
		logChannel := conn.LogChannel()
		for {
			select {
			case logEntry, ok := <-logChannel:
				if !ok {
					log.Println("✓ Log channel closed properly")
					return
				}
				logCount++
				if logCount <= 3 { // Show first 3 logs
					msg := strings.TrimSpace(logEntry.Message())
					if len(msg) > 80 {
						msg = msg[:80] + "..."
					}
					fmt.Printf("LOG #%d [Level %d]: %s\n", logCount, logEntry.Level(), msg)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// 4. Register an event listener for specific events
	conn.RegisterEventListener(eslgo.EventListenAll, func(event *eslgo.Event) {
		eventCount++
		if eventCount <= 3 { // Show first 3 events
			// Check if body is JSON
			bodyStr := string(event.Body)
			isJSON := strings.HasPrefix(strings.TrimSpace(bodyStr), "{")
			fmt.Printf("EVENT #%d: %s (JSON: %v)\n", eventCount, event.GetName(), isJSON)
		}
	})

	log.Println("\n=== TEST RESULTS ===")
	log.Println("Feature 1: JSON Event Parsing - Checking...")
	log.Println("Feature 2: Log Reception - Checking...")
	log.Println("Feature 3: Shutdown Stability - Press Ctrl+C to test...")
	log.Println("")

	// Periodic status update
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r[Status] Events: %d, Logs: %d", eventCount, logCount)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Println("\n\n=== SHUTDOWN TEST ===")
	log.Println("Received shutdown signal, testing clean exit...")
	
	// Cancel context to stop goroutines
	cancel()
	
	// Give goroutines time to exit
	time.Sleep(100 * time.Millisecond)
	
	log.Println("\n=== FINAL RESULTS ===")
	log.Printf("✓ JSON Events Received: %d", eventCount)
	log.Printf("✓ Log Entries Received: %d", logCount)
	log.Println("✓ Clean shutdown completed without panic")
}