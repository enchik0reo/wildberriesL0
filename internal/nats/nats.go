package nats

import (
	"fmt"

	"github.com/enchik0reo/wildberriesL0/pkg/logging"
	"github.com/nats-io/stan.go"
)

type ConnNats interface {
	Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error)
	Close() error
}

type Stan struct {
	conn ConnNats
	log  *logging.Lgr
}

func New(clusterID, clientID, url string, logger *logging.Lgr) (*Stan, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("can't connect to nuts: %w", err)
	}

	return &Stan{conn: sc, log: logger}, nil
}

func (s *Stan) GetMsg(ch chan []byte) {
	_, err := s.conn.Subscribe("orders", func(m *stan.Msg) {
		ch <- m.Data
	}, stan.DeliverAllAvailable())
	if err != nil {
		s.log.Fatalf("can't subscribe: %v\n", err)
	}
}

func (s *Stan) CloseConnect() error {
	return s.conn.Close()
}
