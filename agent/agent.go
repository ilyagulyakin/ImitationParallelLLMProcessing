package agent

import (
	"context"
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
	ErrTemporary    = errors.New("временная ошибка")
	ErrNeedOperator = errors.New("нужна помощь оператора")
)

func (FakeAgent) Process(ctx context.Context, ticket Ticket) error {
	log.Printf("агент принял заявку id=%s text=%q", ticket.ID, ticket.Text)

	delay := time.Duration(rand.Intn(500)) * time.Millisecond

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
	}

	switch n := rand.Intn(100); {
	case n < 80:
		return nil
	case n < 90:
		return ErrTemporary
	case n < 97:
		return ErrNeedOperator
	default:
		return errors.New("агент не смог разобрать ответ модели")
	}
}
