package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	_ "github.com/lib/pq"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=DB
type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Close() error
}

type Storage struct {
	db DB
}

func New(host, port, user, password, dbname, sslmode string) (*Storage, error) {
	connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, o models.Order) error {
	q := `INSERT INTO orders VALUES ($1, $2)`

	if _, err := s.db.ExecContext(ctx, q, o.Uid, o.Details); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetById(uid string) ([]byte, error) {
	q := `SELECT details FROM orders WHERE uid = $1`

	var details []byte

	if err := s.db.QueryRow(q, uid).Scan(&details); err != nil {
		return nil, fmt.Errorf("doesn't exist in db; ")
	}
	return details, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]models.Order, error) {
	q := `SELECT * FROM orders`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("can't load orders from db: %w", err)
	}
	defer rows.Close()

	n, err := rows.Columns()
	if err != nil {
		fmt.Printf("can't get the column names: %v", err)
		n = make([]string, 0)
	}

	var orders = make([]models.Order, 0, len(n))

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Uid, &order.Details)
		if err != nil {
			return nil, fmt.Errorf("can't scan order from db: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) CloseConnect(ctx context.Context) error {
	return s.db.Close()
}
