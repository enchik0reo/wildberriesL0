package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	"github.com/nats-io/stan.go"
)

type Repo interface {
	GetByUid(uid string) ([]byte, error)
	Save(ctx context.Context, order models.Order) error
	Stop(ctx context.Context) error
}

type Stan struct {
	Conn stan.Conn
}

func NatsConnect(clusterID, clientID, url string) (*Stan, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("can't connect to nuts: %w", err)
	}

	return &Stan{Conn: sc}, nil
}

func (s *Stan) GetMsg(ctx context.Context, repo Repo) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	sub, err := s.Conn.Subscribe("orders", func(m *stan.Msg) {
		b := m.Data
		uid, err := validate(b)
		if err != nil {
			log.Printf("can't validate message: %v", err)
			return
		}

		if err = repo.Save(ctx, models.Order{Uid: uid, Details: b}); err != nil {
			log.Printf("can't save message: %v", err)
			return
		}
	}, stan.StartWithLastReceived())
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	wg.Wait()

	if err = sub.Unsubscribe(); err != nil {
		return fmt.Errorf("can't unsubscribe: %w", err)
	}

	return err
}

func validate(message []byte) (string, error) {
	chk := models.Check{}

	err := json.Unmarshal(message, &chk)
	if err != nil {
		return chk.OrderUid, fmt.Errorf("invalid message: %w", err)
	}

	return chk.OrderUid, nil
}
