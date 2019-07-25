package service

import (
	"context"
	"time"

	proto "github.com/mdevilliers/cache-service/proto/v1"

	"github.com/rs/zerolog"
)

type service struct {
	logger zerolog.Logger
	client store
}

type store interface {
	Set(string, string, time.Duration) error
	Get(string) (string, error)
	RandomKey() (string, error)
	Del(string) error
}

func NewCacheService(logger zerolog.Logger, client store) *service {
	return &service{
		logger: logger,
		client: client,
	}
}

func (s *service) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {

	err := s.client.Set(req.GetKey(), req.GetContents(), time.Second*time.Duration(req.GetTtl()))

	if err != nil {

		s.logger.Err(err).Msg("error caching item")

		return &proto.SetResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: err.Error(),
					Code:    proto.ErrorCode_UNKNOWN_ERROR,
				},
			},
		}, err
	}

	return &proto.SetResponse{
		Status: &proto.Status{
			Ok: true,
		},
	}, nil
}
func (s *service) GetByKey(ctx context.Context, req *proto.GetByKeyRequest) (*proto.GetByKeyResponse, error) {

	value, err := s.client.Get(req.GetKey())

	if err != nil {

		s.logger.Err(err).Msg("error getting item")

		return &proto.GetByKeyResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: err.Error(),
					Code:    proto.ErrorCode_UNKNOWN_ERROR,
				},
			},
		}, err
	}

	return &proto.GetByKeyResponse{
		Key:      req.GetKey(),
		Contents: value,
		Status: &proto.Status{
			Ok: true,
		},
	}, nil

}
func (s *service) GetRandomN(ctx context.Context, req *proto.GetRandomNRequest) (*proto.GetRandomNResponse, error) {

	keys := []string{}

	for i := 0; i < int(req.Count); i++ {

		key, err := s.client.RandomKey()

		if err != nil {

			s.logger.Err(err).Msg("error getting random key")

			return &proto.GetRandomNResponse{
				Status: &proto.Status{
					Ok: false,
					Error: &proto.Error{
						Message: err.Error(),
						Code:    proto.ErrorCode_UNKNOWN_ERROR,
					},
				},
			}, err

		}

		keys = append(keys, key)
	}

	return &proto.GetRandomNResponse{
		Keys: keys,
		Status: &proto.Status{
			Ok: true,
		},
	}, nil

}
func (s *service) Purge(ctx context.Context, req *proto.PurgeRequest) (*proto.PurgeResponse, error) {

	err := s.client.Del(req.GetKey())

	if err != nil {

		s.logger.Err(err).Msg("error purging item")

		return &proto.PurgeResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: err.Error(),
					Code:    proto.ErrorCode_UNKNOWN_ERROR,
				},
			},
		}, err
	}

	return &proto.PurgeResponse{
		Status: &proto.Status{
			Ok: true,
		},
	}, nil

}
