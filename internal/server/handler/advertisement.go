package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/database/model"
	"vk-internship/internal/logger"
	"vk-internship/internal/utils"
)

type CreateAdRequest struct {
	Caption     string  `json:"caption" validate:"required,min=3,max=128"`
	Description string  `json:"description" validate:"required,max=1024"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

type CreateAdResponse struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Caption     string    `json:"caption"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url,omitempty"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateAdHandler(cfg *config.ServerConfig, log logger.Logger, db database.Database) http.HandlerFunc {
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
