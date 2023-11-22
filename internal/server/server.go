package server

import (
	"config_center/internal/config"
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	server *http.Server
	router *mux.Router
}

func New(cfg *config.Config, logger *zap.Logger) *Server {
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:              cfg.HttpServerAddress,
		Handler:           r,
		ReadTimeout:       cfg.HttpServerReadTimeout,
		WriteTimeout:      cfg.HttpServerWriteTimeout,
		ReadHeaderTimeout: cfg.HttpServerReadTimeout,
	}

	return &Server{
		cfg:    cfg,
		logger: logger,
		server: srv,
		router: r,
	}
}

func (s *Server) Run(ctx context.Context, errChan chan<- error) {
	idleConnectionsClosed := make(chan struct{})

	go func() {
		<-ctx.Done()

		ctxTimeout, cancel := context.WithTimeout(context.Background(), s.cfg.HttpServerShutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(ctxTimeout); err != nil {
			s.logger.Error("an error occurred during shutdown the HTTP server", zap.Error(err))
		}

		close(idleConnectionsClosed)
	}()

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Error("an error occurred during serving the HTTP server", zap.Error(err))
		errChan <- err
	}

	<-idleConnectionsClosed
}
