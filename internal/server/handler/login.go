package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/logger"
	"vk-internship/internal/utils"
)

// LoginRequest представляет запрос на вход
// @Description Запрос для аутентификации пользователя
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// LoginResponse представляет ответ при успешном входе
// @Description Ответ после успешной аутентификации пользователя
type LoginResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	CurrentUser  *string   `json:"current_user,omitempty"`
	IsAuthorized bool      `json:"is_authorized"`
}

// LoginHandler обрабатывает запросы на вход
// @Summary Аутентификация пользователя
// @Description Проверяет учетные данные пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "Неверный формат запроса или ошибки валидации"
// @Failure 401 {string} string "Неверные учетные данные"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /login [post]
func LoginHandler(cfg *config.ServerConfig, log logger.Logger, db database.Database) http.HandlerFunc {
	validate := utils.NewValidator()

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			currentUserID string
			isAuthorized  bool
		)

		if ctxVal := r.Context().Value("userID"); ctxVal != nil {
			currentUserID = ctxVal.(string)
			isAuthorized = true
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Warnf("invalid request body", map[string]interface{}{"error": err.Error()})
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := validate.Validate(req); err != nil {
			validationErrors := validate.FormatValidationErrors(err)
			log.Warnf("validation failed", map[string]interface{}{"errors": validationErrors})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErrors)
			return
		}

		user, err := db.GetUserByUsername(req.Username)
		if err != nil {
			if errors.Is(err, database.ErrUserNotFound) {
				log.Warnf("user not found", map[string]interface{}{"username": req.Username})
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			log.Error(err, "failed to get user")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			log.Warnf("invalid password", map[string]interface{}{"username": req.Username})
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWTToken(cfg, user.ID, user.Username)
		if err != nil {
			log.Error(err, "failed to generate token")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)

		response := LoginResponse{
			ID:           user.ID,
			Username:     user.Username,
			CreatedAt:    user.CreatedAt,
			IsAuthorized: isAuthorized,
		}

		if isAuthorized {
			response.CurrentUser = &currentUserID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(err, "failed to encode response")
		}

		log.Infof("user logged in", map[string]interface{}{
			"user_id":       user.ID,
			"username":      user.Username,
			"current_user":  currentUserID,
			"is_authorized": isAuthorized,
		})
	}
}
