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
	conn, err := eslgo.Dial("146.190.77.119:8021", "DialB2025!", nil)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer conn.ExitAndClose()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Reload XML
	response, err := conn.SendCommand(ctx, command.API{
		Command: "reloadxml",
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Reload response: %s\n", string(response.Body))

	// Now originate directly using socket
	response, err = conn.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  "{origination_caller_id_number=1001}sofia/internal/1001@146.190.77.119 &socket('146.190.77.119:8084 async full')",
		Background: false,
	})
	if err != nil {
		log.Printf("Error originating: %v", err)
		return
	}

	fmt.Printf("Originate response: %s\n", string(response.Body))
}