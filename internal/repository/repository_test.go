package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/enchik0reo/wildberriesL0/internal/models"
	"github.com/enchik0reo/wildberriesL0/internal/repository/mocks"
)

func TestGetMsg(t *testing.T) {
	type args struct {
		ctx       context.Context
		order     models.Order
		checkBool bool
		errCache  error
		errStor   error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no exists",
			args: args{
				ctx: context.Background(),
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				checkBool: false,
				errCache:  nil,
				errStor:   nil,
			},
			wantErr: false,
		},
		{
			name: "exists everywhere",
			args: args{
				ctx: context.Background(),
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				checkBool: true,
				errCache:  errors.New(""),
				errStor:   errors.New(""),
			},
			wantErr: true,
		},
		{
			name: "cache exists",
			args: args{
				ctx: context.Background(),
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				checkBool: true,
				errCache:  errors.New(""),
				errStor:   nil,
			},
			wantErr: true,
		},
		{
			name: "storage exists",
			args: args{
				ctx: context.Background(),
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				checkBool: true,
				errCache:  nil,
				errStor:   errors.New(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMc := mocks.NewCache(t)
			storageMc := mocks.NewStorage(t)

			cacheMc.
				On("Save", tt.args.order).
				Return(tt.args.errCache)

			storageMc.
				On("Save", tt.args.ctx, tt.args.order).
				Return(tt.args.errStor)

			r := &Repository{
				cache:   cacheMc,
				storage: storageMc,
			}

			err := r.Save(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetByUid(t *testing.T) {
	type args struct {
		order    models.Order
		errCache error
		errStor  error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no exists",
			args: args{
				order: models.Order{
					Uid:     "order1",
					Details: nil,
				},
				errCache: errors.New(""),
				errStor:  errors.New(""),
			},
			wantErr: true,
		},
		{
			name: "exists in storage",
			args: args{
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				errCache: errors.New(""),
				errStor:  nil,
			},
			wantErr: false,
		},
		{
			name: "exists",
			args: args{
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				errCache: nil,
				errStor:  nil,
			},
			wantErr: false,
		},
		{
			name: "exists in cache",
			args: args{
				order: models.Order{
					Uid:     "order1",
					Details: []byte("mocki"),
				},
				errCache: nil,
				errStor:  errors.New(""),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMc := mocks.NewCache(t)
			storageMc := mocks.NewStorage(t)

			cacheMc.
				On("GetById", tt.args.order.Uid).
				Return(tt.args.order.Details, tt.args.errCache)

			storageMc.
				On("GetById", tt.args.order.Uid).
				Return(tt.args.order.Details, tt.args.errStor)

			r := &Repository{
				cache:   cacheMc,
				storage: storageMc,
			}

			details, err := r.GetByUid(tt.args.order.Uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUid() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if ok := equal(details, tt.args.order.Details); !ok {
				t.Errorf("GetByUid() details = %s, want = %s", details, tt.args.order.Details)
				return
			}
		})
	}
}

func equal(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}
