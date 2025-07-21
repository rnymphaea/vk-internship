package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"vk-internship/internal/cache"
	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/logger"
	"vk-internship/internal/server/handler"
	"vk-internship/internal/server/middleware"
)

// @title VK Internship API
// @version 1.0
// @description API для управления объявлениями
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func NewRouter(cfg *config.ServerConfig, log logger.Logger, db database.Database, cache cache.Cache) *chi.Mux {
	router := chi.NewMux()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.LoggingMiddleware(log))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	router.Get("/", handler.Home)
	router.With(middleware.AuthOptionalMiddleware(cfg, log)).Post("/register", handler.RegistrationHandler(cfg, log, db))
	router.With(middleware.AuthOptionalMiddleware(cfg, log)).Post("/login", handler.LoginHandler(cfg, log, db))

	router.With(middleware.AuthOptionalMiddleware(cfg, log)).Get("/ads", handler.GetAdsHandler(log, db))
	router.With(middleware.AuthOptionalMiddleware(cfg, log)).Get("/ads/{id}", handler.GetAdHandler(log, db))

	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthRequiredMiddleware(cfg, log))
		r.Post("/ads", handler.CreateAdHandler(log, db, cache))
		r.Delete("/ads/{id}", handler.DeleteAdHandler(log, db))
		r.Put("/ads/{id}", handler.UpdateAdHandler(log, db))
	})

	return router
}
