package cache

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/enchik0reo/wildberriesL0/internal/models"
)

func TestSave(t *testing.T) {
	type args struct {
		order models.Order
		err   error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exists",
			args: args{
				order: models.Order{
					Uid:     "order1",
					Details: []byte("some details1"),
				},
				err: errors.New(""),
			},
			wantErr: true,
		},
		{
			name: "no exists",
			args: args{
				order: models.Order{
					Uid:     "order2",
					Details: []byte("some details2"),
				},
				err: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				m:       map[string][]byte{"order1": []byte("some details1")},
				slice:   make([]string, 1),
				count:   0,
				size:    1,
				RWMutex: sync.RWMutex{},
			}

			err := c.Save(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCheck(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "exists",
			args: args{
				uid: "order1",
			},
			want: true,
		},
		{
			name: "no exists",
			args: args{
				uid: "order2",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				m:       map[string][]byte{"order1": []byte("some details1")},
				slice:   make([]string, 1),
				count:   0,
				size:    1,
				RWMutex: sync.RWMutex{},
			}

			bl := c.Check(tt.args.uid)
			if bl != tt.want {
				t.Errorf("Check() bool = %v, want = %v", bl, tt.want)
				return
			}
		})
	}
}

func TestGetById(t *testing.T) {
	type args struct {
		uid string
		msg []byte
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exists",
			args: args{
				uid: "order1",
				msg: []byte("some msg1"),
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "no exists",
			args: args{
				uid: "order2",
				msg: nil,
				err: errors.New(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				m:       map[string][]byte{"order1": []byte("some msg1")},
				slice:   make([]string, 1),
				count:   0,
				size:    1,
				RWMutex: sync.RWMutex{},
			}

			msg, err := c.GetById(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(msg, tt.args.msg) {
				t.Errorf("GetById() msg = %s, wantMsg = %s", msg, tt.args.msg)
				return
			}
		})
	}
}
