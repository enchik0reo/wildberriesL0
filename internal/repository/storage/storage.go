package storage

import (
	"context"

	"github.com/enchik0reo/wildberriesL0/internal/models"
)

type Storage interface {
	Save(context.Context, models.Order) error
	GetAll(context.Context) ([]models.Order, error)
	Stop(ctx context.Context) error
}
