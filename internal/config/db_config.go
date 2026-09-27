package config

type DBConfig struct {
	PG_Host     string
	PG_Port     string
	PG_Password string
	PG_User     string
	PG_DBName   string
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		PG_Host:     getEnv("PG_HOST", "localhost"),
		PG_Port:     getEnv("PG_PORT", "5432"),
		PG_Password: getEnv("PG_PASSWORD", "admin"),
		PG_User:     getEnv("PG_USER", "postgres"),
		PG_DBName:   getEnv("PG_DBNAME", "sdn"),
	}
}
