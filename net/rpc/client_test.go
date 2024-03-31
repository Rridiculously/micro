package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitClientProxy(t *testing.T) {

	testcases := []struct {
		name string

		service Service
		wantErr error
	}{{
		name:    "nil",
		service: nil,
		wantErr: errors.New("service is nil"),
	},
		{
			name:    "no pointer",
			service: &UserService{},
			wantErr: errors.New("service is not a pointer"),
		},
		{
			name:    "UserService",
			service: &UserService{},
			wantErr: errors.New("service is not a pointer"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service, "")
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.(*UserService).GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			t.Log(resp)

		})
	}
}
