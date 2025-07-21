package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command"
	"github.com/yosefkohan26/eslgo/command/call"
)

// Handle outbound connections
func handleOutbound(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	log.Println("=== OUTBOUND CONNECTION RECEIVED ===")
	
	// Test 1: GetVar
	log.Println("\nTEST 1: Testing GetVar command...")
	varValue, err := conn.GetVar(ctx, "my_test_variable")
	if err != nil {
		log.Printf("❌ GetVar failed: %v", err)
	} else {
		log.Printf("✅ GetVar success! Value: '%s'", varValue)
	}

	// Test 2: Set a variable using call.Set
	log.Println("\nTEST 2: Testing Set variable...")
	uuid := response.ChannelUUID()
	_, err = conn.SendCommand(ctx, call.Set{
		UUID:  uuid,
		Key:   "variable_from_go",
		Value: "Hello from eslgo!",
	})
	if err != nil {
		log.Printf("❌ Set variable failed: %v", err)
	} else {
		log.Printf("✅ Set variable success!")
	}

	// Test 3: Resume
	log.Println("\nTEST 3: Testing Resume command...")
	resumeResp, err := conn.Resume(ctx)
	if err != nil {
		log.Printf("❌ Resume failed: %v", err)
	} else {
		log.Printf("✅ Resume success! Response: %v", resumeResp.IsOk())
	}
	
	log.Println("\n=== TESTS COMPLETED ===")
}

func main() {
	// Start outbound listener in goroutine
	go func() {
		log.Println("Starting outbound ESL server on :8084...")
		if err := eslgo.ListenAndServe(":8084", handleOutbound); err != nil {
			log.Fatalf("Outbound server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(2 * time.Second)

	// Connect to FreeSWITCH
	log.Println("Connecting to FreeSWITCH...")
	conn, err := eslgo.Dial("146.190.77.119:8021", "DialB2025!", nil)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer conn.ExitAndClose()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a simple inline dialplan that sets a variable and connects to our socket
	dialplanApp := fmt.Sprintf("set:my_test_variable=TestValue123,socket:146.190.77.119:8084 async full")
	
	log.Println("Originating test call with inline dialplan...")
	response, err := conn.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  fmt.Sprintf("{origination_caller_id_number=1001}sofia/internal/1001@146.190.77.119 '%s' inline", dialplanApp),
		Background: false,
	})
	if err != nil {
		log.Printf("Error originating: %v", err)
		return
	}

	log.Printf("Originate response: %s", string(response.Body))
	
	// Wait for tests to complete
	time.Sleep(5 * time.Second)
	log.Println("Test completed!")
}