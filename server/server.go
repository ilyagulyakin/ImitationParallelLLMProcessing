package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"ImitationParallelLLMProcessing/agent"
)

// Server принимает NDJSON-поток тикетов и передаёт их в канал.
//
// Канал создаёт и закрывает main. Server только пишет в него. Закрывать канал
// можно только после того, как Shutdown вернул nil (активных хендлеров не осталось).
type Server struct {
	ch       chan<- agent.Ticket
	srv      *http.Server
	accepted atomic.Int64
}

func NewServer(addr string, ch chan<- agent.Ticket) *Server {
	mux := http.NewServeMux()

	s := &Server{
		ch:  ch,
		srv: &http.Server{Addr: addr, Handler: mux},
	}

	mux.HandleFunc("/tickets", s.ticketHandler)
	return s
}

// Run блокируется, пока сервер не остановят через Shutdown.
// Штатная остановка (http.ErrServerClosed) не считается ошибкой.
func (s *Server) Run() error {
	log.Printf("SERVER: слушаю %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown перестаёт принимать новые соединения и ждёт завершения активных
// хендлеров (включая дочитывание тела). nil означает: хендлеров не осталось,
// канал можно закрывать.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// Accepted возвращает число тикетов, переданных в канал.
func (s *Server) Accepted() int64 {
	return s.accepted.Load()
}

func (s *Server) ticketHandler(w http.ResponseWriter, r *http.Request) {
	tickets, err := decodeTickets(r)
	if err != nil {
		log.Printf("SERVER: не смог разобрать тело: %v", err)
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	accepted := 0
	for _, t := range tickets {
		select {
		case s.ch <- t:
			accepted++
			s.accepted.Add(1)
		case <-r.Context().Done():
			// Shutdown не отменяет r.Context(), поэтому сюда попадаем только
			// если отвалился сам клиент.
			log.Printf("SERVER: клиент отвалился, принято %d из %d", accepted, len(tickets))
			return
		}
	}

	log.Printf("SERVER: принято тикетов: %d", accepted)
	return
}

func decodeTickets(r *http.Request) ([]agent.Ticket, error) {
	dec := json.NewDecoder(r.Body)

	var tickets []agent.Ticket
	for {
		var t agent.Ticket
		err := dec.Decode(&t)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if t.ID == "" {
			return nil, errors.New("ticket.id is required")
		}
		tickets = append(tickets, t)
	}
	if len(tickets) == 0 {
		return nil, errors.New("empty body")
	}
	return tickets, nil
}
