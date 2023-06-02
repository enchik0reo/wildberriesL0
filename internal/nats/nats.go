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
	sub  stan.Subscription
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
	sub, err := s.conn.Subscribe("orders", func(m *stan.Msg) {
		ch <- m.Data
	}, stan.DeliverAllAvailable())
	if err != nil {
		s.log.Fatalf("can't subscribe: %v\n", err)
	}
	s.sub = sub
}

func (s *Stan) CloseConnect() error {
	var errors []byte

	if err := s.sub.Unsubscribe(); err != nil {
		e := fmt.Sprintf("can't unsubscribe to the channel: %s; ", err.Error())
		errors = append(errors, []byte(e)...)
	}

	if err := s.conn.Close(); err != nil {
		e := fmt.Sprintf("can't close connection to the cluster.: %s; ", err.Error())
		errors = append(errors, []byte(e)...)
	}

	if len(errors) != 0 {
		return fmt.Errorf("%s", errors)
	}

	return nil
}
