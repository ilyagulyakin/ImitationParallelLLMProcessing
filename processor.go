package sever

import (
	"ImitationParallelLLMProcessing/agent"
	"context"
	"log"
)

func Run(ctx context.Context, ch chan<- agent.Ticket) {
	for {
		select {
		case <-ctx.Done():
			return // Продумать логирование
		case ticket, ok := <-ch:
			if !ok {
				return // Продумать логирование
			}

			log.Printf("обрабатываю %s: %q", ticket.ID, ticket.Text)
			// здесь можно вызвать agent.FakeAgent{}.Process(...)
		}

	}
}
