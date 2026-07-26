package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	appauth "lebaotai.com/backend-go/internal/app/auth"
	"lebaotai.com/backend-go/internal/infra/persistence/memory"
	"lebaotai.com/backend-go/internal/infra/security"
	intauth "lebaotai.com/backend-go/internal/interface/auth"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	secret := os.Getenv("JWT_SECRET")
	ttl := os.Getenv("JWT_TTL")
	d, e := time.ParseDuration(ttl)
	if e != nil {
		log.Fatalf("invalid JWT_TTL: %v", e)
	}

	// userRepo := postgres.NewPostgresUserRepository(pool)
	userRepo := memory.NewInMemoryUserRepository()
	hasher := security.NewBcryptHasher()
	jwtIssuer := security.NewJWTIssuer(secret, d)

	// create usecase handlers
	// registerUseCase := auth.NewRegisterUserHandler(userRepo, hasher)
	loginUseCase := appauth.NewLoginUserHandler(userRepo, hasher, jwtIssuer)

	// create HTTP handlers
	authHTTPHandler := intauth.NewAuthHandler(
		loginUseCase,
	)

	r := gin.Default()

	// auth routes
	r.POST("/auth/login", authHTTPHandler.Login)

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("======== Backend is starting with port 8080 ========")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
