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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First set the dialplan variable manually
	response, err := conn.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  "{origination_caller_id_number=1001,my_dialplan_variable=FreeSWITCH_is_awesome}sofia/internal/1001@146.190.77.119 &socket('192.168.4.22:8084 async full')",
		Background: false,
	})
	if err != nil {
		log.Printf("Error originating: %v", err)
		return
	}

	fmt.Printf("Originate response: %s\n", string(response.Body))
}