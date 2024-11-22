package transport

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
	"time"
)

type HTTPServerConfig interface {
	GetAddress() string
}

type HTTPServer struct {
	log    zerolog.Logger
	server *http.Server
	wg     sync.WaitGroup
}

func NewHTTPServer(log zerolog.Logger, router http.Handler, cfg HTTPServerConfig) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    cfg.GetAddress(),
			Handler: router,
		},
		log: log.With().Str("module", "http-server").Logger(),
	}
}

func (s *HTTPServer) MustStart() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.log.Info().Str("address", s.server.Addr).Msg("Starting HTTP server")

		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal().Err(err).Msg("HTTP server start failed")
		}

		s.log.Info().Msg("HTTP server stopped")
	}()
}

func (s *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.log.Info().Msg("Stopping HTTP server")
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Fatal().Err(err).Msg("HTTP server shutdown failed")
	}

	s.wg.Wait()
	s.log.Info().Msg("HTTP server shutdown complete")
}
