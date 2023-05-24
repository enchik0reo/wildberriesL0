package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/enchik0reo/wildberriesL0/internal/models"
	"github.com/enchik0reo/wildberriesL0/internal/repository/storage"
)

type Cacher interface {
	Save(o models.Order)
	GetById(uid string) ([]byte, error)
	Chk(uid string) ([]byte, bool)
}

type Repository struct {
	storage storage.Storage
	cache   Cacher
}

func New(s storage.Storage, c Cacher) *Repository {
	orders, err := s.GetAll(context.Background())
	if err != nil {
		log.Printf("can't warmup cache: %s\n", err.Error())
	}

	for _, o := range orders {
		c.Save(o)
	}

	return &Repository{
		storage: s,
		cache:   c,
	}
}

func (r *Repository) Save(ctx context.Context, order models.Order) error {
	if _, ok := r.cache.Chk(order.Uid); !ok {
		if err := r.storage.Save(ctx, order); err != nil {
			return err
		}
		r.cache.Save(order)
	} else {
		log.Printf("order with uid: %s exists", order.Uid)
	}
	return nil
}

func (r *Repository) GetByUid(uid string) ([]byte, error) {
	details, err := r.cache.GetById(uid)
	if err != nil {
		details, err = r.storage.GetById(uid)
		if err != nil {
			return nil, fmt.Errorf("order with uid: %s doesn't exist", uid)
		}
	}
	return details, nil
}

func (s *Repository) Stop(ctx context.Context) error {
	return s.storage.Stop(ctx)
}
