package service

import (
	"context"
	"errors"
	"testing"

	"github.com/enchik0reo/wildberriesL0/internal/service/mocks"
)

func TestStop(t *testing.T) {
	type args struct {
		ctx     context.Context
		errNats error
		errRepo error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:     context.Background(),
				errNats: nil,
				errRepo: nil,
			},
			wantErr: false,
		},
		{
			name: "nats err",
			args: args{
				ctx:     context.Background(),
				errNats: errors.New(""),
				errRepo: nil,
			},
			wantErr: true,
		},
		{
			name: "repo err",
			args: args{
				ctx:     context.Background(),
				errNats: nil,
				errRepo: errors.New(""),
			},
			wantErr: true,
		},
		{
			name: "all err",
			args: args{
				ctx:     context.Background(),
				errNats: errors.New(""),
				errRepo: errors.New(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			natsMc := mocks.NewNats(t)
			repoMc := mocks.NewRepository(t)

			natsMc.
				On("CloseConnect").
				Return(tt.args.errNats)

			repoMc.
				On("CloseConnect", tt.args.ctx).
				Return(tt.args.errRepo)

			s := &Service{
				nats: natsMc,
				repo: repoMc,
			}

			err := s.Stop(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}
