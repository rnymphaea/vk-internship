package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/database/model"
	"vk-internship/internal/logger"
	"vk-internship/internal/utils"
)

// RegistrationRequest представляет запрос на регистрацию
// @Description Запрос для регистрации нового пользователя
type RegistrationRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// RegistrationResponse представляет ответ при успешной регистрации
// @Description Ответ после успешной регистрации пользователя
type RegistrationResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	Token        string    `json:"token"`
	CurrentUser  *string   `json:"current_user,omitempty"`
	IsAuthorized bool      `json:"is_authorized"`
}

// RegistrationHandler обрабатывает запросы на регистрацию
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegistrationRequest true "Данные для регистрации"
// @Success 201 {object} RegistrationResponse
// @Failure 400 {object} map[string]string "Неверный формат запроса или ошибки валидации"
// @Failure 409 {string} string "Пользователь с таким именем уже существует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /register [post]
func RegistrationHandler(cfg *config.ServerConfig, log logger.Logger, db database.Database) http.HandlerFunc {
	validate := utils.NewValidator()

	return func(w http.ResponseWriter, r *http.Request) {
		var currentUserID string
		var isAuthorized bool

		if ctxVal := r.Context().Value("userID"); ctxVal != nil {
			currentUserID = ctxVal.(string)
			isAuthorized = true
		}

		var req RegistrationRequest
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err, "password hashing failed")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user := &model.User{
			Username: req.Username,
			Password: string(hashedPassword),
		}

		createdUser, err := db.CreateUser(user)
		if err != nil {
			if errors.Is(err, database.ErrUserExists) {
				log.Warnf("username already taken", map[string]interface{}{"username": req.Username})
				http.Error(w, "Username already exists", http.StatusConflict)
				return
			}

			log.Error(err, "failed to create user")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		token, err := utils.GenerateJWTToken(cfg, createdUser.ID, createdUser.Username)
		if err != nil {
			log.Error(err, "failed to generate token")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)

		response := RegistrationResponse{
			ID:           createdUser.ID,
			Username:     createdUser.Username,
			CreatedAt:    createdUser.CreatedAt,
			Token:        token,
			IsAuthorized: isAuthorized,
		}

		if isAuthorized {
			response.CurrentUser = &currentUserID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(err, "failed to encode response")
		}

		log.Infof("user registered", map[string]interface{}{
			"user_id":       createdUser.ID,
			"username":      createdUser.Username,
			"created_at":    createdUser.CreatedAt,
			"token":         token,
			"is_authorized": isAuthorized,
		})
	}
}
