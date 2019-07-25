package service

import (
	"context"
	"testing"

	"github.com/mdevilliers/cache-service/internal/service/mocks"
	proto "github.com/mdevilliers/cache-service/proto/v1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test_Cache_OK(t *testing.T) {

	fake := &mocks.FakeStore{}
	log := zerolog.New(nil).With().Logger()
	ctx := context.Background()

	service := NewCacheService(log, fake)

	response, err := service.Set(ctx, &proto.SetRequest{Key: "foo", Contents: "bar"})

	require.Nil(t, err)
	require.True(t, response.GetStatus().Ok)

}

func Test_Cache_Errors(t *testing.T) {

	fake := &mocks.FakeStore{}

	storeErr := errors.New("boom")
	fake.SetReturns(storeErr)

	log := zerolog.New(nil).With().Logger()
	ctx := context.Background()

	service := NewCacheService(log, fake)

	response, err := service.Set(ctx, &proto.SetRequest{Key: "foo", Contents: "bar"})

	require.NotNil(t, err)
	require.Equal(t, storeErr, errors.Cause(err))

	require.False(t, response.GetStatus().Ok)
}
