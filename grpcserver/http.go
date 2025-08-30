package grpcserver

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type HTTPGateway struct {
	logger   *zap.Logger
	server   *http.Server
	grpcPort string
	httpPort string
	mux      *runtime.ServeMux
	opts     []grpc.DialOption
}

func NewHTTPGateway(logger *zap.Logger, grpcPort, httpPort string, opts ...grpc.DialOption) *HTTPGateway {
	return &HTTPGateway{
		logger:   logger,
		mux:      runtime.NewServeMux(),
		opts:     opts,
		grpcPort: grpcPort,
		httpPort: httpPort,
	}
}

func (g *HTTPGateway) RegisterService(register func(mux *runtime.ServeMux, host string, opts []grpc.DialOption) error) error {
	return register(g.mux, fmt.Sprintf("localhost:%s", g.grpcPort), g.opts)
}

func (g *HTTPGateway) Start() error {
	g.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", g.httpPort),
		Handler: g.mux,
	}
	g.logger.Info("starting HTTP server", zap.String("port", g.httpPort))

	return g.server.ListenAndServe()
}

func (g *HTTPGateway) Stop() error {
	if g.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return g.server.Shutdown(ctx)
}
