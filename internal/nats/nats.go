package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	"github.com/nats-io/stan.go"
)

type Repository interface {
	Save(ctx context.Context, order models.Order) error
}

type Stan struct {
	conn stan.Conn
	repo Repository
}

func New(clusterID, clientID, url string, rp Repository) (*Stan, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("can't connect to nuts: %w", err)
	}

	return &Stan{conn: sc, repo: rp}, nil
}

func (s *Stan) GetMsg(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	sub, err := s.conn.Subscribe("orders", func(m *stan.Msg) {
		b := m.Data
		uid, err := validate(b)
		if err != nil {
			log.Printf("can't validate message: %v", err)
			return
		}

		if err = s.repo.Save(ctx, models.Order{Uid: uid, Details: b}); err != nil {
			log.Printf("can't save message: %v", err)
			return
		}
	}, stan.DeliverAllAvailable())
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	wg.Wait()

	if err = sub.Unsubscribe(); err != nil {
		return fmt.Errorf("can't unsubscribe: %w", err)
	}

	return err
}

func (s *Stan) CloseConnect() error {
	return s.conn.Close()
}

func validate(message []byte) (string, error) {
	chk := models.Basic{}

	err := json.Unmarshal(message, &chk)
	if err != nil {
		return chk.OrderUid, fmt.Errorf("json format: %w", err)
	}

	if chk.OrderUid == "" {
		return chk.OrderUid, fmt.Errorf("invalid message")
	}

	return chk.OrderUid, nil
}
