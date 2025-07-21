package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"vk-internship/internal/cache"
	"vk-internship/internal/database"
	"vk-internship/internal/database/model"
	"vk-internship/internal/logger"
	"vk-internship/internal/utils"
)

// CreateAdRequest представляет запрос на создание объявления
// @Description Данные для создания нового объявления
type CreateAdRequest struct {
	Caption     string  `json:"caption" validate:"required,min=3,max=128"`
	Description string  `json:"description" validate:"required,max=1024"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

// CreateAdResponse представляет ответ после создания объявления
// @Description Информация о созданном объявлении
type CreateAdResponse struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Caption     string    `json:"caption"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url,omitempty"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateAdHandler создает новое объявление
// @Security BearerAuth
// @Summary Создать объявление
// @Description Создает новое объявление от имени авторизованного пользователя
// @Tags ads
// @Accept json
// @Produce json
// @Param request body CreateAdRequest true "Данные объявления"
// @Success 201 {object} CreateAdResponse
// @Failure 400 {object} map[string]string "Неверный формат запроса или ошибки валидации"
// @Failure 401 {string} string "Не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /ads [post]
func CreateAdHandler(log logger.Logger, db database.Database, cache cache.Cache) http.HandlerFunc {
	validate := utils.NewValidator()

	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			log.Warn("userID not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateAdRequest
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

		ad := &model.Advertisement{
			AuthorID:    userID,
			Caption:     req.Caption,
			Description: req.Description,
			ImageURL:    req.ImageURL,
			Price:       int(req.Price * 100),
		}

		createdAd, err := db.CreateAd(ad)
		if err != nil {
			log.Error(err, "failed to create ad")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		go func() {
			if err := cache.UpdateFeed(context.TODO(), *createdAd); err != nil {
				log.Warn("failed to update feed cache")
			}
		}()

		response := CreateAdResponse{
			ID:          createdAd.ID,
			AuthorID:    createdAd.AuthorID,
			Caption:     createdAd.Caption,
			Description: createdAd.Description,
			ImageURL:    createdAd.ImageURL,
			Price:       float64(createdAd.Price) / 100,
			CreatedAt:   createdAd.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(err, "failed to encode response")
		}

		log.Infof("advertisement created", map[string]interface{}{
			"advertisement_id": createdAd.ID,
			"author_id":        createdAd.AuthorID,
			"caption":          createdAd.Caption,
			"description":      createdAd.Description,
			"image_url":        createdAd.ImageURL,
			"price":            createdAd.Price,
			"created_at":       createdAd.CreatedAt,
		})
	}
}

// GetAdResponse представляет информацию об объявлении
// @Description Полная информация об объявлении
type GetAdResponse struct {
	ID             string    `json:"id"`
	AuthorUsername string    `json:"author_username"`
	Caption        string    `json:"caption"`
	Description    string    `json:"description"`
	ImageURL       string    `json:"image_url,omitempty"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	IsOwner        *bool     `json:"is_owner,omitempty"`
}

// GetAdHandler возвращает информацию об объявлении
// @Summary Получить объявление
// @Description Возвращает полную информацию об объявлении по ID
// @Tags ads
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Security ApiKeyAuth
// @Success 200 {object} GetAdResponse
// @Failure 400 {string} string "Неверный ID объявления"
// @Failure 404 {string} string "Объявление не найдено"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /ads/{id} [get]
func GetAdHandler(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adID := chi.URLParam(r, "id")
		if adID == "" {
			log.Warn("id not provided")
			http.Error(w, "Ad ID is required", http.StatusBadRequest)
			return
		}

		ad, err := db.GetAd(r.Context(), adID)
		if err != nil {
			if errors.Is(err, database.ErrAdNotFound) {
				log.Warnf("ad not found", map[string]interface{}{"ad_id": adID})
				http.Error(w, "Ad not found", http.StatusNotFound)
				return
			}
			log.Error(err, "failed to get ad")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var userID string
		var isAuthenticated bool
		if ctxUserID, ok := r.Context().Value("userID").(string); ok && ctxUserID != "" {
			userID = ctxUserID
			isAuthenticated = true
		}

		response := GetAdResponse{
			ID:             ad.ID,
			AuthorUsername: ad.AuthorUsername,
			Caption:        ad.Caption,
			Description:    ad.Description,
			ImageURL:       ad.ImageURL,
			Price:          float64(ad.Price) / 100,
			CreatedAt:      ad.CreatedAt,
		}

		if isAuthenticated {
			isOwner := userID == ad.AuthorID
			response.IsOwner = &isOwner
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(err, "failed to encode response")
		}
	}
}

// DeleteAdHandler удаляет объявление
// @Security BearerAuth
// @Summary Удалить объявление
// @Description Удаляет объявление по ID (только для автора объявления)
// @Tags ads
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Security BearerAuth
// @Success 204 "Объявление успешно удалено"
// @Failure 400 {string} string "Неверный ID объявления"
// @Failure 401 {string} string "Не авторизован"
// @Failure 403 {string} string "Нет прав на удаление"
// @Failure 404 {string} string "Объявление не найдено"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /ads/{id} [delete]
func DeleteAdHandler(log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		adID := chi.URLParam(r, "id")
		if adID == "" {
			http.Error(w, "Ad ID is required", http.StatusBadRequest)
			return
		}

		err := db.DeleteAd(r.Context(), adID, userID)
		if err != nil {
			if errors.Is(err, database.ErrAdNotFoundOrNotOwnedByUser) {
				http.Error(w, "Ad not found or not owned by user", http.StatusNotFound)
				return
			}
			log.Error(err, "failed to delete ad")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// UpdateAdRequest представляет запрос на обновление объявления
// @Description Данные для обновления объявления (все поля опциональны)
type UpdateAdRequest struct {
	Caption     string  `json:"caption" validate:"omitempty,min=3,max=128"`
	Description string  `json:"description" validate:"omitempty,max=1024"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
}

// UpdateAdResponse представляет ответ после обновления объявления
// @Description Информация об обновленном объявлении
type UpdateAdResponse struct {
	ID          string    `json:"id"`
	Caption     string    `json:"caption"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url,omitempty"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateAdHandler обновляет объявление
// @Security ApiKeyAuth
// @Summary Обновить объявление
// @Description Обновляет информацию об объявлении (только для автора объявления)
// @Tags ads
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Param request body UpdateAdRequest true "Данные для обновления"
// @Security BearerAuth
// @Success 200 {object} UpdateAdResponse
// @Failure 400 {object} map[string]string "Неверный формат запроса или ошибки валидации"
// @Failure 401 {string} string "Не авторизован"
// @Failure 403 {string} string "Нет прав на обновление"
// @Failure 404 {string} string "Объявление не найдено"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /ads/{id} [put]
func UpdateAdHandler(log logger.Logger, db database.Database) http.HandlerFunc {
	validate := utils.NewValidator()

	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		adID := chi.URLParam(r, "id")
		if adID == "" {
			http.Error(w, "Ad ID is required", http.StatusBadRequest)
			return
		}

		var req UpdateAdRequest
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
		}

		currentAd, err := db.GetAd(r.Context(), adID)
		if err != nil {
			if errors.Is(err, database.ErrAdNotFound) {
				http.Error(w, "Ad not found", http.StatusNotFound)
				return
			}
			log.Error(err, "failed to get ad")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if currentAd.AuthorID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		update := model.Advertisement{
			ID:       adID,
			AuthorID: userID,
		}

		if req.Caption != "" {
			update.Caption = req.Caption
		} else {
			update.Caption = currentAd.Caption
		}

		if req.Description != "" {
			update.Description = req.Description
		} else {
			update.Description = currentAd.Description
		}

		if req.ImageURL != "" {
			update.ImageURL = req.ImageURL
		} else {
			update.ImageURL = currentAd.ImageURL
		}

		if req.Price > 0 {
			update.Price = int(req.Price * 100)
		} else {
			update.Price = currentAd.Price
		}

		updatedAd, err := db.UpdateAd(r.Context(), &update)
		if err != nil {
			log.Error(err, "failed to update ad")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := UpdateAdResponse{
			ID:          updatedAd.ID,
			Caption:     updatedAd.Caption,
			Description: updatedAd.Description,
			ImageURL:    updatedAd.ImageURL,
			Price:       float64(updatedAd.Price) / 100,
			CreatedAt:   updatedAd.CreatedAt,
			UpdatedAt:   updatedAd.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(err, "failed to encode response")
		}
	}
}
