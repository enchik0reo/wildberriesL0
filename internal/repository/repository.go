package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	_ "github.com/vektra/mockery"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Storage
type Storage interface {
	Save(context.Context, models.Order) error
	GetById(uid string) ([]byte, error)
	GetAll(context.Context) ([]models.Order, error)
	CloseConnect(ctx context.Context) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Cache
type Cache interface {
	Save(o models.Order) error
	GetById(uid string) ([]byte, error)
	Check(uid string) bool
}

type Repository struct {
	storage Storage
	cache   Cache
}

func New(s Storage, c Cache) *Repository {
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
	var errors []byte
	if err := r.cache.Save(order); err != nil {
		e := fmt.Sprintf("in cache: %s; ", err.Error())
		errors = append(errors, []byte(e)...)
	}

	if err := r.storage.Save(ctx, order); err != nil {
		e := fmt.Sprintf("in db: %s; ", err.Error())
		errors = append(errors, []byte(e)...)
	}

	if len(errors) != 0 {
		return fmt.Errorf("%s", errors)
	}

	return nil
}

func (r *Repository) GetByUid(uid string) ([]byte, error) {
	details, err := r.cache.GetById(uid)
	if err == nil {
		return details, err
	}

	details, err = r.storage.GetById(uid)
	if err == nil {
		r.cache.Save(models.Order{Uid: uid, Details: details})
		return details, err
	}

	return nil, fmt.Errorf("doesn't exist")
}

func (s *Repository) CloseConnect(ctx context.Context) error {
	return s.storage.CloseConnect(ctx)
}
