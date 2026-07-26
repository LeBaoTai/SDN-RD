package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDBConnection(ctx context.Context) (*pgxpool.Pool, error) {
	log.Println("Initializing pool...")

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	userDBUser := os.Getenv("DB_USER")
	userDBPassword := os.Getenv("DB_PASSWORD")
	userDBName := os.Getenv("DB_NAME")
	userDBPort := os.Getenv("DB_PORT")
	userDBHost := os.Getenv("DB_HOST")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", userDBUser, userDBPassword, userDBHost, userDBPort, userDBName)

	connConfig, err := pgxpool.ParseConfig(connString)

	// setting connection pool para
	connConfig.MaxConns = 20
	connConfig.MinConns = 5
	connConfig.MaxConnLifetime = 30 * time.Minute
	connConfig.MaxConnIdleTime = 5 * time.Minute
	connConfig.HealthCheckPeriod = 1 * time.Minute

	// connect to the database
	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
