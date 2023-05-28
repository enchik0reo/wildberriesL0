package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/enchik0reo/wildberriesL0/internal/models"
)

type Nats interface {
	GetMsg(chan []byte)
	CloseConnect() error
}

type Repository interface {
	Save(context.Context, models.Order) error
	CloseConnect(ctx context.Context) error
}

type Service struct {
	nats Nats
	repo Repository
}

func New(nats Nats, repo Repository) *Service {
	return &Service{
		nats: nats,
		repo: repo,
	}
}

func (s *Service) Work(ctx context.Context) {
	ch := make(chan []byte)

	go s.nats.GetMsg(ch)

	for {
		msg := <-ch

		uid, err := validate(msg)
		if err != nil {
			log.Printf("can't validate message [%s] error: %v", msg, err)
			continue
		}

		if err = s.repo.Save(ctx, models.Order{Uid: uid, Details: msg}); err != nil {
			log.Printf("can't save message with uid [%s] error: %v", uid, err)
			continue
		}
	}
}

func (s *Service) Stop(ctx context.Context) error {
	var errors []byte

	if err := s.nats.CloseConnect(); err != nil {
		s := fmt.Sprintf("can't close nats connection: %s; ", err.Error())
		errors = append(errors, []byte(s)...)
	}

	if err := s.repo.CloseConnect(ctx); err != nil {
		s := fmt.Sprintf("can't close repository connection: %s; ", err.Error())
		errors = append(errors, []byte(s)...)
	}

	if len(errors) != 0 {
		return fmt.Errorf("can't stop service: %s", errors)
	}

	return nil
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
