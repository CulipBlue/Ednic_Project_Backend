package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv                    string
	AppPort                   string
	FrontendURL               string
	DatabaseDSN               string
	DatabaseMaxOpenConns      int
	DatabaseMaxIdleConns      int
	DatabaseConnMaxLifetimeMi int
	JWTSecret                 string
	JWTAccessTokenTTLMinutes  int
	SuperAdminName            string
	SuperAdminUsername        string
	SuperAdminEmail           string
	SuperAdminPassword        string
}

func Load() Config {
	loadDotEnv(".env")

	return Config{
		AppEnv:                    getEnv("APP_ENV", "local"),
		AppPort:                   getEnv("APP_PORT", "8080"),
		FrontendURL:               getEnv("FRONTEND_URL", "http://localhost:3000"),
		DatabaseDSN:               getEnv("DATABASE_DSN", ""),
		DatabaseMaxOpenConns:      getEnvInt("DATABASE_MAX_OPEN_CONNS", 25),
		DatabaseMaxIdleConns:      getEnvInt("DATABASE_MAX_IDLE_CONNS", 25),
		DatabaseConnMaxLifetimeMi: getEnvInt("DATABASE_CONN_MAX_LIFETIME_MINUTES", 5),
		JWTSecret:                 getEnv("JWT_SECRET", "local-development-secret"),
		JWTAccessTokenTTLMinutes:  getEnvInt("JWT_ACCESS_TOKEN_TTL_MINUTES", 60),
		SuperAdminName:            getEnv("SUPER_ADMIN_NAME", ""),
		SuperAdminUsername:        getEnv("SUPER_ADMIN_USERNAME", ""),
		SuperAdminEmail:           getEnv("SUPER_ADMIN_EMAIL", ""),
		SuperAdminPassword:        getEnv("SUPER_ADMIN_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
