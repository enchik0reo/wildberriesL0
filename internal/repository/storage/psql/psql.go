package psql

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/enchik0reo/wildberriesL0/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
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

	tableScript := readScript()

	if _, err = db.Exec(tableScript); err != nil {
		return nil, fmt.Errorf("can't read table script: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, o models.Order) error {
	q := `INSERT INTO orders VALUES ($1, $2)`

	if _, err := s.db.ExecContext(ctx, q, o.Uid, o.Details); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) GetById(uid string) ([]byte, error) {
	q := `SELECT details FROM orders WHERE uid = $1`

	row := s.db.QueryRow(q, uid)

	var details []byte

	if err := row.Scan(&details); err != nil {
		return nil, fmt.Errorf("order with uid: %s doesn't exist in db", uid)
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

func (s *Storage) Stop(ctx context.Context) error {
	return s.db.Close()
}

func readScript() string {
	fname := "script/init_table.txt"

	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("can't open script file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("can't read script file: %v", err)
	}

	return string(bytes)
}
