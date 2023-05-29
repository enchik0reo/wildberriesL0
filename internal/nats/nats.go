package nats

import (
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
)

type NatsConn interface {
	Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error)
	Close() error
}

type Stan struct {
	conn NatsConn
}

func New(clusterID, clientID, url string) (*Stan, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("can't connect to nuts: %w", err)
	}

	return &Stan{conn: sc}, nil
}

func (s *Stan) GetMsg(ch chan []byte) {
	_, err := s.conn.Subscribe("orders", func(m *stan.Msg) {
		ch <- m.Data
	}, stan.DeliverAllAvailable())
	if err != nil {
		log.Fatalf("can't subscribe: %v", err)
	}
}

func (s *Stan) CloseConnect() error {
	return s.conn.Close()
}
