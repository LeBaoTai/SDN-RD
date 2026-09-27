package router

import (
	"github.com/LeBaoTai/SDN-RD/internal/backend/handler"
	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	IntentHandler *handler.IntentHandler
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.HealthCheck)

	v1 := r.Group("/api")
	{
		RegisterIntentRoutes(v1, deps.IntentHandler)
	}

	return r
}
