package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/enchik0reo/wildberriesL0/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Nats
type Nats interface {
	GetMsg(chan []byte)
	CloseConnect() error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Repository
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
			log.Printf("can't save order with uid [%s]: %v", uid, err)
			continue
		}
	}
}

func (s *Service) Stop(ctx context.Context) error {
	var errors []byte

	if err := s.nats.CloseConnect(); err != nil {
		e := fmt.Sprintf("can't close nats connection: %s; ", err.Error())
		errors = append(errors, []byte(e)...)
	}

	if err := s.repo.CloseConnect(ctx); err != nil {
		e := fmt.Sprintf("can't close repository connection: %s; ", err.Error())
		errors = append(errors, []byte(e)...)
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
