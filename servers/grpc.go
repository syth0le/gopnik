package servers

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

const (
	defaultServerName = "default server"
)

type GRPCServerConfig struct {
	Port             int64 `yaml:"port"`
	EnableRecover    bool  `yaml:"enable_recover"` // TODO write recover and tls
	EnableReflection bool  `yaml:"enable_reflection"`
}

func (c *GRPCServerConfig) Validate() error {
	return nil
}

type grpcServerOptions struct {
	serverName               string
	serverOptions            []grpc.ServerOption
	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor
	//healthCheckHandler TODO
}

type GRPCServerOption func(serverOptions *grpcServerOptions)

func GRPCWithServerName(serverName string) GRPCServerOption {
	return func(serverOptions *grpcServerOptions) {
		serverOptions.serverName = serverName
	}
}

func GRPCWithUnaryInterceptors(ins ...grpc.UnaryServerInterceptor) GRPCServerOption {
	return func(serverOptions *grpcServerOptions) {
		serverOptions.unaryServerInterceptors = ins
	}
}

func GRPCWithStreamInterceptors(ins ...grpc.StreamServerInterceptor) GRPCServerOption {
	return func(serverOptions *grpcServerOptions) {
		serverOptions.streamServerInterceptors = ins
	}
}

func GRPCWithServerOptions(opts ...grpc.ServerOption) GRPCServerOption {
	return func(serverOptions *grpcServerOptions) {
		serverOptions.serverOptions = opts
	}
}

type GRPCServer struct {
	Server     *grpc.Server
	logger     *zap.Logger
	serverName string
	address    string
}

func NewGRPCServer(cfg GRPCServerConfig, logger *zap.Logger, opts ...GRPCServerOption) *GRPCServer {
	serverOptions := &grpcServerOptions{
		serverName:               defaultServerName,
		serverOptions:            nil,
		unaryServerInterceptors:  nil,
		streamServerInterceptors: nil,
	}

	for _, opt := range opts {
		opt(serverOptions)
	}

	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(serverOptions.unaryServerInterceptors...),
		grpc.ChainStreamInterceptor(serverOptions.streamServerInterceptors...),
	}
	grpcOptions = append(grpcOptions, serverOptions.serverOptions...)

	server := grpc.NewServer(grpcOptions...)

	if cfg.EnableReflection {
		reflection.Register(server)
	}

	return &GRPCServer{
		Server:     server,
		logger:     logger,
		serverName: serverOptions.serverName,
		address:    fmt.Sprintf(":%d", cfg.Port),
	}
}

func (s *GRPCServer) Run() error {
	s.logger.Sugar().Infof("starting grpc server %s at %s", s.serverName, s.address)

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("cannot listen grpc server %s on %s: %w", s.serverName, s.address, err)
	}

	err = s.Server.Serve(l)
	if err != nil {
		return fmt.Errorf("serve grpc server %s on %s: %w", s.serverName, s.address, err)
	}

	return nil
}

func (s *GRPCServer) GracefullyStop() error {
	s.logger.Sugar().Infof("gracefully stopping grpc server %s on: %s", s.serverName, s.address)
	s.Server.GracefulStop()
	return nil
}

func (s *GRPCServer) ForcefullyStop() error {
	s.logger.Sugar().Infof("forcefully stopping grpc server %s on: %s", s.serverName, s.address)
	s.Server.Stop()
	return nil
}
