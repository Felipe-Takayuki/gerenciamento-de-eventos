package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	CORSOrigins []string
}

func Load() *Config {
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	var origins []string
	for _, origin := range strings.Split(corsOrigins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	return &Config{
		Port:        getEnv("PORT", "3000"),
		DBHost:      getEnv("DB_HOST", "127.0.0.1"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "adamas_db"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-jwt-key"),
		CORSOrigins: origins,
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
