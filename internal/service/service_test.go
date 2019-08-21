package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mdevilliers/cache-service/internal/service/mocks"
	"github.com/mdevilliers/cache-service/internal/store"
	proto "github.com/mdevilliers/cache-service/proto/v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test_Cache(t *testing.T) {

	fake := &mocks.FakeStore{}
	log := zerolog.New(nil).With().Logger()
	ctx := context.Background()

	service := NewCacheService(log, fake)

	testCases := []struct {
		name     string
		request  *proto.SetRequest
		errCode  proto.ErrorCode
		err      error
		storeErr error
	}{
		{
			name:    "ok",
			request: &proto.SetRequest{Key: "foo", Contents: "bar"},
		},
		{
			name:    "no key",
			request: &proto.SetRequest{Contents: "bar"},
			errCode: proto.ErrorCode_KEY_NOT_SUPPLIED,
			err:     ErrNoKeySupplied,
		},
		{
			name:    "no contents",
			request: &proto.SetRequest{Key: "foo"},
			errCode: proto.ErrorCode_CONTENT_NOT_SUPPLIED,
			err:     ErrNoContentSupplied,
		},
		{
			name:     "store returns an error",
			request:  &proto.SetRequest{Key: "foo", Contents: "bar"},
			errCode:  proto.ErrorCode_UNKNOWN_ERROR,
			err:      errors.New("boom"),
			storeErr: errors.New("boom"),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			fake.SetReturns(tc.storeErr)

			response, err := service.Set(ctx, tc.request)

			if tc.err == nil {
				require.Nil(t, err)
				require.True(t, response.GetStatus().Ok)

			} else {
				require.False(t, response.GetStatus().Ok)
				require.Equal(t, tc.err.Error(), response.GetStatus().GetError().GetMessage())
				require.Equal(t, tc.errCode, response.GetStatus().GetError().GetCode())
			}
		},
		)
	}
}

func Test_GetByKey(t *testing.T) {

	fake := &mocks.FakeStore{}
	log := zerolog.New(nil).With().Logger()
	ctx := context.Background()

	service := NewCacheService(log, fake)

	testCases := []struct {
		name     string
		request  *proto.GetByKeyRequest
		errCode  proto.ErrorCode
		err      error
		storeErr error
		statusOK bool
	}{
		{
			name:     "ok",
			request:  &proto.GetByKeyRequest{Key: "foo"},
			statusOK: true,
		},
		{
			name:    "no key",
			request: &proto.GetByKeyRequest{},
			errCode: proto.ErrorCode_KEY_NOT_SUPPLIED,
			err:     ErrNoKeySupplied,
		},
		{
			name:     "store returns an error",
			request:  &proto.GetByKeyRequest{Key: "foo"},
			errCode:  proto.ErrorCode_UNKNOWN_ERROR,
			err:      errors.New("boom"),
			storeErr: errors.New("boom"),
		},
		{
			name:     "item not found",
			request:  &proto.GetByKeyRequest{Key: "foo"},
			errCode:  proto.ErrorCode_KEY_NOT_FOUND,
			statusOK: true, // NOT FOUND is not an error
			storeErr: store.ErrItemNotFound,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			fake.GetReturns("bar", tc.storeErr)

			response, err := service.GetByKey(ctx, tc.request)

			if tc.err == nil {
				require.Nil(t, err)

			} else {
				require.Equal(t, tc.err.Error(), response.GetStatus().GetError().GetMessage())
				require.Equal(t, tc.errCode, response.GetStatus().GetError().GetCode())
			}
			require.Equal(t, tc.statusOK, response.GetStatus().GetOk())
		},
		)
	}
}

func Test_GetRandomN(t *testing.T) {

	fake := &mocks.FakeStore{}
	log := zerolog.New(nil).With().Logger()
	ctx := context.Background()

	service := NewCacheService(log, fake)

	testCases := []struct {
		name     string
		request  *proto.GetRandomNRequest
		errCode  proto.ErrorCode
		err      error
		storeErr error
	}{
		{
			name:    "ok",
			request: &proto.GetRandomNRequest{Count: 5},
		},
		{
			name:    "no count",
			request: &proto.GetRandomNRequest{},
			errCode: proto.ErrorCode_COUNT_NOT_SUPPLIED,
			err:     ErrNoCountSupplied,
		},
		{
			name:     "store returns an error",
			request:  &proto.GetRandomNRequest{Count: 5},
			errCode:  proto.ErrorCode_UNKNOWN_ERROR,
			err:      errors.New("boom"),
			storeErr: errors.New("boom"),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			fake.RandomKeyReturns("bar", tc.storeErr)

			response, err := service.GetRandomN(ctx, tc.request)

			if tc.err == nil {
				require.Nil(t, err)
				require.True(t, response.GetStatus().Ok)

			} else {
				require.False(t, response.GetStatus().Ok)
				require.Equal(t, tc.err.Error(), response.GetStatus().GetError().GetMessage())
				require.Equal(t, tc.errCode, response.GetStatus().GetError().GetCode())

			}
		},
		)
	}
}

func Test_Purge(t *testing.T) {

	fake := &mocks.FakeStore{}
	log := zerolog.New(nil).With().Logger()
	ctx := context.Background()

	service := NewCacheService(log, fake)

	testCases := []struct {
		name     string
		request  *proto.PurgeRequest
		errCode  proto.ErrorCode
		err      error
		storeErr error
	}{
		{
			name:    "ok",
			request: &proto.PurgeRequest{Key: "foo"},
		},
		{
			name:    "no key",
			request: &proto.PurgeRequest{},
			errCode: proto.ErrorCode_KEY_NOT_SUPPLIED,
			err:     ErrNoKeySupplied,
		},
		{
			name:     "store returns an error",
			request:  &proto.PurgeRequest{Key: "foo"},
			errCode:  proto.ErrorCode_UNKNOWN_ERROR,
			err:      errors.New("boom"),
			storeErr: errors.New("boom"),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			fake.DelReturns(tc.storeErr)

			response, err := service.Purge(ctx, tc.request)

			if tc.err == nil {
				require.Nil(t, err)
				require.True(t, response.GetStatus().Ok)

			} else {
				require.False(t, response.GetStatus().Ok)
				require.Equal(t, tc.err.Error(), response.GetStatus().GetError().GetMessage())
				require.Equal(t, tc.errCode, response.GetStatus().GetError().GetCode())

			}
		},
		)
	}
}
