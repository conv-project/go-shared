package grpcserver

import (
	"fmt"
	"google.golang.org/grpc/reflection"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server represents gRPC server.
type Server struct {
	server *grpc.Server
	port   string
	logger *zap.Logger
}

// New creates a new gRPC server.
func New(logger *zap.Logger, port string, opts ...grpc.ServerOption) *Server {
	server := grpc.NewServer(opts...)

	return &Server{
		server: server,
		port:   port,
		logger: logger,
	}
}

// RegisterService registers a service with the server.
func (s *Server) RegisterService(register func(server *grpc.Server)) {
	register(s.server)
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	reflection.Register(s.server)

	addr := fmt.Sprintf(":%s", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.logger.Info("Starting gRPC server", zap.String("port", s.port))

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop stops the gRPC server.
func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		s.logger.Info("gRPC server stopped")
	}
}
