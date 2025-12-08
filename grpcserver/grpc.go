package grpcserver

import (
	"fmt"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server represents gRPC server.
type Server struct {
	server *grpc.Server
	health *health.Server
	port   string
}

// New creates a new gRPC server.
func New(port string, opts ...grpc.ServerOption) *Server {
	server := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)

	return &Server{
		server: server,
		health: healthServer,
		port:   port,
	}
}

// RegisterService registers a service with the server.
func (s *Server) RegisterService(register func(server *grpc.Server)) {
	register(s.server)
}

func (s *Server) SetServingStatus(name string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.health.SetServingStatus(name, status)
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	reflection.Register(s.server)

	addr := fmt.Sprintf(":%s", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	slog.Info("starting gRPC server", slog.String("port", s.port))

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop stops the gRPC server.
func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		slog.Info("gRPC server stopped")
	}
}
