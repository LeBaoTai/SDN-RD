package router

import (
	"github.com/LeBaoTai/SDN-RD/internal/backend/handler"
	"github.com/gin-gonic/gin"
)

func RegisterIntentRoutes(rg *gin.RouterGroup, h *handler.IntentHandler) {
	intents := rg.Group("/intents")
	{
		intents.POST("", h.CreateIntent)
	}
}
