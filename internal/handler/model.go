package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/006lp/akashchat-api-go/internal/model"
	"github.com/gin-gonic/gin"
)

// ModelHandler handles model-related HTTP requests
type ModelHandler struct{}

// NewModelHandler creates a new ModelHandler instance
func NewModelHandler() *ModelHandler {
	return &ModelHandler{}
}

// GetModels handles the /v1/models endpoint
func (h *ModelHandler) GetModels(c *gin.Context) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://chat.akash.network/api/models/", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Referer", "https://chat.akash.network/")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch models"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var models []model.Model
	if err := json.Unmarshal(body, &models); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal models"})
		return
	}

	c.JSON(http.StatusOK, models)
}