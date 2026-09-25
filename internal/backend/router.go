package backend

import (
	"github.com/LeBaoTai/SDN-RD/internal/backend/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(intentHandler *handler.IntentHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.HealthCheck)
	r.POST("/intents", intentHandler.CreateIntent)

	return r
}
