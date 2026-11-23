package main

import (
	"ImitationParallelLLMProcessing/agent"
	"log"
	"sync"
	"sync/atomic"
)

var processed atomic.Int64

// Processed возвращает число завершённых тикетов.
func Processed() int64 {
	return processed.Load()
}

// Run обрабатывает тикеты, пока канал не закрыт. Выход только по close(ch).
// Всё, что успело попасть в канал, будет обработано.
func Run(ch <-chan agent.Ticket, workers int) {
	var wg sync.WaitGroup
	wg.Add(workers)

	for id := 1; id <= workers; id++ {
		go func(id int) {
			defer wg.Done()
			for ticket := range ch {
				log.Printf("PROCESSOR[%d]: обрабатываю %s: %s", id, ticket.ID, ticket.Text)
				process(id, ticket)
				processed.Add(1)
			}
		}(id)
	}

	wg.Wait()
	log.Println("PROCESSOR: все воркеры завершены")
}

func process(worker int, ticket agent.Ticket) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PROCESSOR[%d]: PANIC на тикете %s: %v", worker, ticket.ID, r)
		}
	}()

	if err := (agent.FakeAgent{}).Process(ticket); err != nil {
		log.Printf("PROCESSOR[%d]: ticket %s: %v", worker, ticket.ID, err)
	}
}
