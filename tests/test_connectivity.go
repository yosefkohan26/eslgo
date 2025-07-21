package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command"
)

func main() {
	// FreeSWITCH server details from the Python script
	host := "146.190.77.119:8021"
	password := "DialB2025!"

	log.Printf("Connecting to FreeSWITCH at %s...", host)

	conn, err := eslgo.Dial(host, password, func() {
		log.Println("Connection disconnected.")
	})
	if err != nil {
		log.Fatalf("Error connecting to FreeSWITCH: %v", err)
	}
	defer conn.ExitAndClose()

	log.Println("✓ Successfully connected to FreeSWITCH")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test basic command - get FreeSWITCH status
	response, err := conn.SendCommand(ctx, command.API{Command: "status"})
	if err != nil {
		log.Printf("Error getting status: %v", err)
	} else {
		fmt.Println("\nFreeSWITCH Status:")
		fmt.Println("-" + string(make([]byte, 50)))
		fmt.Println(string(response.Body))
	}

	// Check gateway status
	response, err = conn.SendCommand(ctx, command.API{
		Command:   "sofia",
		Arguments: "status gateway sip_voipessential",
	})
	if err != nil {
		log.Printf("Error getting gateway status: %v", err)
	} else {
		fmt.Println("\nGateway Status:")
		fmt.Println("-" + string(make([]byte, 50)))
		fmt.Println(string(response.Body))
	}

	log.Println("\n✓ Connectivity test completed successfully")
}