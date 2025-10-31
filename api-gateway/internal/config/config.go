package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppPort         int
	AppEnv          string
	UserServiceAddr string
	CardServiceAddr string
	LogLevel        string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		AppName:         getEnv("APP_NAME", "api-gateway"),
		AppPort:         getEnvInt("APP_PORT", 8080),
		AppEnv:          getEnv("APP_ENV", "development"),
		UserServiceAddr: fmt.Sprintf("%s:%s", getEnv("USER_SERVICE_HOST", "localhost"), getEnv("USER_SERVICE_PORT", "50051")),
		CardServiceAddr: fmt.Sprintf("%s:%s", getEnv("CARD_SERVICE_HOST", "localhost"), getEnv("CARD_SERVICE_PORT", "50052")),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

