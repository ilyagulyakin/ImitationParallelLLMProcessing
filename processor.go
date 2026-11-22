package main

import (
	"ImitationParallelLLMProcessing/agent"
	"context"
	"log"
)

func Run(ctx context.Context, ch <-chan agent.Ticket) {
	for {
		select {
		case <-ctx.Done():
			log.Println("процессор: контекст отменён, выхожу")
			return
		case ticket, ok := <-ch:
			if !ok {
				log.Println("процессор: канал закрыт, выхожу")
				return
			}

			// лог ДО обработки — «взял в работу»
			log.Printf("обрабатываю %s: %s", ticket.ID, ticket.Text)

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("PANIC на тикете %s: %v", ticket.ID, r)
					}
				}()
				if err := (agent.FakeAgent{}).Process(ctx, ticket); err != nil {
					log.Printf("ticket %s: %v", ticket.ID, err)
				}
			}()
		}
	}
}
