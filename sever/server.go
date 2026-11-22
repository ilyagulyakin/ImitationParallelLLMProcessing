package server

import (
	"ImitationParallelLLMProcessing/agent"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type Server struct {
	ch chan<- agent.Ticket
}

func NewServer(ch chan<- agent.Ticket) *Server {
	return &Server{ch: ch}
}

func (s *Server) ticketHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Println("ticketHandler получил запрос")

	decoder := json.NewDecoder(r.Body)

	for {
		var ticket agent.Ticket
		if err := decoder.Decode(&ticket); errors.Is(err, io.EOF) {
			break

		} else if err != nil {
			log.Println("не смог декодировать JSON:", err)
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}
		s.ch <- ticket
	}

}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/tickets", s.ticketHandler)
	return http.ListenAndServe(addr, mux)
}

//
//func ticketHandler(w http.ResponseWriter, r *http.Request, ch chan agent.Ticket) {
//	log.Println("ticketHandler получил JSON")
//	decoder := json.NewDecoder(r.Body)
//	defer r.Body.Close()
//
//	for {
//		var ticket agent.Ticket
//		if err := decoder.Decode(&ticket); errors.Is(err, io.EOF) {
//			break
//		} else if err != nil {
//			log.Println("Сервер не смог декодировать JSON:", err)
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//	}
//
//}
//
//func runserver() {
//	http.HandleFunc("/tickets", ticketHandler)
//	log.Fatal(http.ListenAndServe(":9091", nil))
//}
