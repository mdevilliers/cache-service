package server

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	proto "github.com/mdevilliers/cache-service/proto/v1"
)

type server struct {
	logger     zerolog.Logger
	grpcServer *grpc.Server
	service    proto.CacheServer
}

// New returns a configured GRPC service
func New(logger zerolog.Logger, service proto.CacheServer) *server {
	return &server{
		logger:  logger,
		service: service,
	}
}

// Start successfully starts the server or returns an error
// Usually this server definition would be in a library and
// configure middleware, tracing and metrics
func (s *server) Start(address string) error {

	// Create a new listener on the defined address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error creating listener")
	}

	serv := grpc.NewServer(
		// clients should reconnect regulary so we can loadbalace effectively
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
	)
	s.grpcServer = serv

	proto.RegisterCacheServer(serv, s.service)

	s.logger.Info().Fields(map[string]interface{}{
		"address": address,
	}).Msg("starting server")

	if err := serv.Serve(lis); err != nil {
		return errors.Wrap(err, "error running server")
	}
	return nil
}

// Stop attempts a graceful shutdown of th server
func (s *server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
