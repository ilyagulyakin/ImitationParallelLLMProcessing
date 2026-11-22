package main

import (
	"context"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	//f, err := os.Open("tickets.json")
	//if err != nil {
	//	return err
	//}
	//defer f.Close()
	//
	//dec := json.NewDecoder(f)
	//for {
	//	var ticket agent.Ticket
	//	if err := dec.Decode(&ticket); errors.Is(err, io.EOF) {
	//		break
	//	} else if err != nil {
	//		return err
	//	}
	//
	//	if err := (agent.FakeAgent{}).Process(ctx, ticket); err != nil {
	//		log.Printf("ticket %s: %v", ticket.ID, err)
	//	}
	//}
	//
	//return nil
}
