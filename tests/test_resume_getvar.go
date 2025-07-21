package main

import (
	"context"
	"log"

	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command/call"
)

func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	log.Println("New connection received from FreeSWITCH")

	// 1. Get the variable set by the dialplan
	initialVar, err := conn.GetVar(ctx, "my_dialplan_variable")
	if err != nil {
		log.Printf("Error getting variable: %v", err)
		return
	}
	log.Printf("Successfully received variable from dialplan: %s", initialVar)

	// 2. Set a new variable on the channel
	uuid := response.ChannelUUID()
	_, err = conn.SendCommand(ctx, call.Set{
		UUID:  uuid,
		Key:   "variable_from_go",
		Value: "Hello from eslgo!",
	})
	if err != nil {
		log.Printf("Error setting variable: %v", err)
		return
	}
	log.Println("Successfully set a new variable on the channel")

	// 3. Resume dialplan execution
	log.Println("Resuming dialplan execution...")
	_, err = conn.Resume(ctx)
	if err != nil {
		log.Printf("Error resuming dialplan: %v", err)
	}
}

func main() {
	log.Println("Starting outbound ESL server on :8084")
	log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
}