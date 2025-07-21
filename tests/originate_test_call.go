package main

import (
	"context"
	"log"
	"time"

	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command"
)

func main() {
	// Connect to FreeSWITCH
	conn, err := eslgo.Dial("146.190.77.119:8021", "DialB2025!", func() {
		log.Println("Connection disconnected.")
	})
	if err != nil {
		log.Fatalf("Error connecting to FreeSWITCH: %v", err)
	}
	defer conn.ExitAndClose()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Originate a call to extension 5000
	log.Println("Originating call to extension 5000...")
	response, err := conn.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  "user/1000 5000 XML default",
		Background: false,
	})
	if err != nil {
		log.Printf("Error originating call: %v", err)
		return
	}

	log.Printf("Originate response: %s", string(response.Body))
	
	// Wait a bit to see the results
	time.Sleep(10 * time.Second)
}