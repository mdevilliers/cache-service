module github.com/mdevilliers/cache-service

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/maxbrunsfeld/counterfeiter/v6 v6.4.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.0.0 // indirect
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v0.0.7
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.39.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
)
