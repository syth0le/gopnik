package clients

import (
	"context"
	"fmt"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	defaultClientName = "default client"
)

type GRPCClientConnConfig struct {
	Endpoint  string `yaml:"endpoint"`
	UserAgent string `yaml:"user_agent"`

	MaxRetries            int64         `yaml:"max_retries"`
	TimeoutBetweenRetries time.Duration `yaml:"timeout_between_retries"`

	InitTimeout time.Duration `yaml:"init_timeout"`

	EnableCompressor bool `yaml:"enable_compressor"`
}

func (c *GRPCClientConnConfig) Validate() error {
	return nil
}

type grpcClientConnOptions struct {
	clientName  string
	dialOptions []grpc.DialOption
}

type GRPCClientConnOption func(serverOptions *grpcClientConnOptions)

func GRPCWithClientName(clientName string) GRPCClientConnOption {
	return func(clientOptions *grpcClientConnOptions) {
		clientOptions.clientName = clientName
	}
}

func GRPCWithDialOptions(dialOpts ...grpc.DialOption) GRPCClientConnOption {
	return func(serverOptions *grpcClientConnOptions) {
		serverOptions.dialOptions = dialOpts
	}
}

func NewGRPCClientConn(ctx context.Context, cfg GRPCClientConnConfig, opts ...GRPCClientConnOption) (*grpc.ClientConn, error) {
	clientConnOptions := &grpcClientConnOptions{
		clientName:  defaultClientName,
		dialOptions: nil,
	}

	for _, opt := range opts {
		opt(clientConnOptions)
	}

	dialOptions := []grpc.DialOption{
		grpc.WithUserAgent(cfg.UserAgent),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(
				retry.WithMax(uint(cfg.MaxRetries)),
				retry.WithPerRetryTimeout(cfg.TimeoutBetweenRetries),
			),
		),
	}

	if cfg.EnableCompressor {
		dialOptions = append(dialOptions, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}

	ctx, cancelFunc := context.WithTimeout(ctx, cfg.InitTimeout)
	defer cancelFunc()

	conn, err := grpc.DialContext(ctx, cfg.Endpoint, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial context on %s: %w", cfg.Endpoint, err)
	}

	return conn, nil
}
