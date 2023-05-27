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
}

type Repository interface {
	Save(context.Context, models.Order) error
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
			log.Printf("can't validate message: [%s] error: %v", msg, err)
			continue
		}

		if err = s.repo.Save(ctx, models.Order{Uid: uid, Details: msg}); err != nil {
			log.Printf("can't save message: [%s] error: %v", msg, err)
			continue
		}
	}
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
