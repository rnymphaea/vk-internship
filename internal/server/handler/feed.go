package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/logger"
)

type FeedResponse struct {
	Ads        []AdResponse `json:"ads"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	Total      int          `json:"total"`
	TotalPages int          `json:"total_pages"`
}

type AdResponse struct {
	ID             string    `json:"id"`
	AuthorUsername string    `json:"author_username"`
	Caption        string    `json:"caption"`
	Description    string    `json:"description"`
	ImageURL       string    `json:"image_url,omitempty"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	IsOwner        *bool     `json:"is_owner,omitempty"`
}

var ValidSorts = map[string]struct{}{
	"created_at": {},
	"price":      {},
}

func GetAdsHandler(cfg *config.ServerConfig, log logger.Logger, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(query.Get("page_size"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}

		sortBy := query.Get("sort_by")
		if _, ok := ValidSorts[sortBy]; !ok {
			sortBy = "created_at"
		}

		order := strings.ToUpper(query.Get("order"))
		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		var minPrice, maxPrice *int
		if minStr := query.Get("min_price"); minStr != "" {
			if val, err := strconv.ParseFloat(minStr, 64); err == nil && val >= 0 {
				minPriceVal := int(val * 100)
				minPrice = &minPriceVal
			}
		}

		if maxStr := query.Get("max_price"); maxStr != "" {
			if val, err := strconv.ParseFloat(maxStr, 64); err == nil && val >= 0 {
				maxPriceVal := int(val * 100)
				maxPrice = &maxPriceVal
			}
		}

		if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
			log.Error(err, "min_price > max_price")
			http.Error(w, "min_price must be less than or equal to max_price", http.StatusBadRequest)
			return
		}

		ads, total, err := db.GetAds(
			r.Context(),
			sortBy,
			order,
			minPrice,
			maxPrice,
			page,
			pageSize,
		)
		if err != nil {
			log.Error(err, "failed to get ads")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		totalPages := total / pageSize
		if total%pageSize > 0 {
			totalPages++
		}

		if totalPages > 0 && page > totalPages {
			page = totalPages
			ads, total, err = db.GetAds(
				r.Context(),
				sortBy,
				order,
				minPrice,
				maxPrice,
				page,
				pageSize,
			)
			if err != nil {
				log.Error(err, "failed to get ads")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		var userID string
		var isAuthenticated bool
		if ctxUserID, ok := r.Context().Value("userID").(string); ok && ctxUserID != "" {
			userID = ctxUserID
			isAuthenticated = true
		}

		log.Debugf("check userID", map[string]interface{}{"userID": userID, "isAuthenticated": isAuthenticated})

		responseAds := make([]AdResponse, 0, len(ads))
		for _, ad := range ads {
			respAd := AdResponse{
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
				respAd.IsOwner = &isOwner
			}

			responseAds = append(responseAds, respAd)
		}

		response := FeedResponse{
			Ads:        responseAds,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(err, "failed to encode response")
		}
	}
}
