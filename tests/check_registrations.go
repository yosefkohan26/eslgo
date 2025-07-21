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

	// Check registered users
	response, err := conn.SendCommand(ctx, command.API{
		Command:   "sofia",
		Arguments: "status profile internal reg",
	})
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Registered users:")
	fmt.Println(string(response.Body))

	// Also check dialplan
	response, err = conn.SendCommand(ctx, command.API{
		Command:   "xml_locate",
		Arguments: "dialplan default destination_number 5000",
	})
	if err != nil {
		log.Printf("Error checking dialplan: %v", err)
	} else {
		fmt.Println("\nDialplan for 5000:")
		fmt.Println(string(response.Body))
	}
}