package psql

import (
	"context"
	"errors"
	"testing"

	"github.com/enchik0reo/wildberriesL0/internal/models"
	"github.com/enchik0reo/wildberriesL0/internal/repository/storage/psql/mocks"
)

func TestSave(t *testing.T) {
	type args struct {
		ctx     context.Context
		order   models.Order
		query   string
		execErr error
		rez     interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				order: models.Order{
					Uid:     "order1",
					Details: []byte("details..."),
				},
				query:   `INSERT INTO orders VALUES ($1, $2)`,
				execErr: nil,
			},
			wantErr: false,
		},
		{
			name: "exec error",
			args: args{
				ctx: context.Background(),
				order: models.Order{
					Uid:     "order1",
					Details: []byte("details..."),
				},
				query:   `INSERT INTO orders VALUES ($1, $2)`,
				execErr: errors.New(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMc := mocks.NewDB(t)

			dbMc.
				On("ExecContext", tt.args.ctx, tt.args.query, tt.args.order.Uid, tt.args.order.Details).
				Return(tt.args.rez, tt.args.execErr)

			s := &Storage{
				db: dbMc,
			}

			err := s.Save(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}
