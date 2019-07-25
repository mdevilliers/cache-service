package service

import (
	"context"
	"errors"
	"time"

	proto "github.com/mdevilliers/cache-service/proto/v1"

	"github.com/rs/zerolog"
)

type service struct {
	logger zerolog.Logger
	client store
}

// store allows the mocking of the underlying client
type store interface {
	Set(string, string, time.Duration) error
	Get(string) (string, error)
	RandomKey() (string, error)
	Del(string) error
}

var (
	ErrNoKeySupplied     = errors.New("no key supplied")
	ErrNoContentSupplied = errors.New("no content supplied")
	ErrNoCountSupplied   = errors.New("no count supplied")
)

func NewCacheService(logger zerolog.Logger, client store) *service {
	return &service{
		logger: logger,
		client: client,
	}
}

func (s *service) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {

	if req.GetKey() == "" {
		return &proto.SetResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: ErrNoKeySupplied.Error(),
					Code:    proto.ErrorCode_KEY_NOT_SUPPLIED,
				},
			}}, nil
	}

	if req.GetContents() == "" {
		return &proto.SetResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: ErrNoContentSupplied.Error(),
					Code:    proto.ErrorCode_CONTENT_NOT_SUPPLIED,
				},
			}}, nil
	}

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
		}, nil
	}

	return &proto.SetResponse{
		Status: &proto.Status{
			Ok: true,
		},
	}, nil
}
func (s *service) GetByKey(ctx context.Context, req *proto.GetByKeyRequest) (*proto.GetByKeyResponse, error) {

	if req.GetKey() == "" {
		return &proto.GetByKeyResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: ErrNoKeySupplied.Error(),
					Code:    proto.ErrorCode_KEY_NOT_SUPPLIED,
				},
			}}, nil
	}

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
		}, nil
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

	if req.GetCount() == 0 {
		return &proto.GetRandomNResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: ErrNoCountSupplied.Error(),
					Code:    proto.ErrorCode_COUNT_NOT_SUPPLIED,
				},
			}}, nil
	}

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
			}, nil

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

	if req.GetKey() == "" {
		return &proto.PurgeResponse{
			Status: &proto.Status{
				Ok: false,
				Error: &proto.Error{
					Message: ErrNoKeySupplied.Error(),
					Code:    proto.ErrorCode_KEY_NOT_SUPPLIED,
				},
			}}, nil
	}

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
		}, nil
	}

	return &proto.PurgeResponse{
		Status: &proto.Status{
			Ok: true,
		},
	}, nil

}
