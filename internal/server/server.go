package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"vk-internship/internal/config"
	"vk-internship/internal/logger"
)

type Server struct {
	server *http.Server
	log    logger.Logger
}

func New(cfg *config.ServerConfig, router *chi.Mux, log logger.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        router,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
			MaxHeaderBytes: 1 << 20,
		},
		log: log,
	}
}

func (s *Server) Start() error {
	s.log.Debugf("starting HTTP server",
		map[string]interface{}{
			"port":          s.server.Addr,
			"read_timeout":  s.server.ReadTimeout.String(),
			"write_timeout": s.server.WriteTimeout.String(),
			"idle_timeout":  s.server.IdleTimeout.String(),
		})

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("shutting down server")
	return s.server.Shutdown(ctx)
}
