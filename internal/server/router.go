package server

import (
	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "vk-internship/internal/cache"
	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/logger"
	"vk-internship/internal/server/handler"
	"vk-internship/internal/server/middleware"
)

func NewRouter(cfg *config.ServerConfig, log logger.Logger, db database.Database) *chi.Mux {
	router := chi.NewMux()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.LoggingMiddleware(log))

	router.Get("/", handler.Home)
	router.Post("/register", handler.RegistrationHandler(cfg, log, db))
	router.Post("/login", handler.LoginHandler(cfg, log, db))

	router.With(middleware.AuthOptionalMiddleware(cfg, log)).Get("/ads", handler.GetAdsHandler(log, db))
	router.With(middleware.AuthOptionalMiddleware(cfg, log)).Get("/ads/{id}", handler.GetAdHandler(log, db))

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthRequiredMiddleware(cfg, log))
		r.Post("/ads", handler.CreateAdHandler(cfg, log, db))
		r.Delete("/ads/{id}", handler.DeleteAdHandler(log, db))
		r.Put("/ads/{id}", handler.UpdateAdHandler(log, db))
	})

	return router
}
