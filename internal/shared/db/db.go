package db

import (
	"strings"

	"github.com/LeBaoTai/SDN-RD/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg config.DBConfig) (*gorm.DB, error) {
	var dsn strings.Builder
	dsn.WriteString("host=")
	dsn.WriteString(cfg.PG_Host)
	dsn.WriteString(" ")

	dsn.WriteString("user=")
	dsn.WriteString(cfg.PG_User)
	dsn.WriteString(" ")

	dsn.WriteString("password=")
	dsn.WriteString(cfg.PG_Password)
	dsn.WriteString(" ")

	dsn.WriteString("dbname=")
	dsn.WriteString(cfg.PG_DBName)
	dsn.WriteString(" ")

	dsn.WriteString("port=")
	dsn.WriteString(cfg.PG_Port)
	dsn.WriteString(" ")

	dsn.WriteString("sslmode=disable ")
	dsn.WriteString("TimeZone=Asia/Ho_Chi_Minh")

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
