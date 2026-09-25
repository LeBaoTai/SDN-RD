package handler

import (
	"net/http"

	"github.com/LeBaoTai/SDN-RD/internal/backend/model"
	"github.com/LeBaoTai/SDN-RD/internal/backend/service"
	"github.com/gin-gonic/gin"
)

type IntentHandler struct {
	service *service.IntentService
}

func NewIntentHandler(s *service.IntentService) *IntentHandler {
	return &IntentHandler{service: s}
}

func (h *IntentHandler) CreateIntent(c *gin.Context) {
	var intent model.Intent
	if err := c.ShouldBindJSON(&intent); err != nil {
		c.JSON(http.StatusBadRequest, model.IntentResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if err := h.service.HandleIntent(intent); err != nil {
		c.JSON(http.StatusInternalServerError, model.IntentResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, model.IntentResponse{
		Status:  "accepted",
		Message: "intent published to controller",
	})
}
