package agent

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type Ticket struct {
	ID   string    `json:"id"`
	Text string    `json:"text"`
	Date time.Time `json:"date"`
}

type FakeAgent struct{}

var (
	ErrTemporary    = errors.New("AGENT: временная ошибка")
	ErrNeedOperator = errors.New("AGENT: нужна помощь оператора")
)

func (FakeAgent) Process(ticket Ticket) error {
	log.Printf("AGENT: принял заявку id=%s text=%q", ticket.ID, ticket.Text)

	delay := time.Duration(rand.Intn(5)) * time.Second

	time.Sleep(delay)

	switch n := rand.Intn(100); {
	case n < 80:
		return nil
	case n < 90:
		return ErrTemporary
	case n < 97:
		return ErrNeedOperator
	default:
		panic("AGENT: не смог разобрать ответ модели")
	}
}
