package middleware

import (
	"context"
	"net/http"
	"strings"

	"vk-internship/internal/config"
	"vk-internship/internal/logger"
	"vk-internship/internal/utils"
)

func AuthMiddleware(cfg *config.ServerConfig, log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("missing authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Warn("invalid authorization format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := utils.VerifyJWTToken(cfg, token)
			if err != nil {
				log.Warnf("invalid token", map[string]interface{}{"error": err.Error()})
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
