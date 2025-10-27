package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	RateLimiter RateLimiterConfig
}
type ServerConfig struct {
	GRPCPort string
	Host     string
}
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  int // in minutes
	RefreshTokenTTL int // in minutes
}
type RateLimiterConfig struct {
	RequestsPerSecond int
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/user-service/")

	v.AutomaticEnv()

	v.BindEnv("server.grpc_port", "GRPC_PORT")
	v.BindEnv("server.host", "SERVER_HOST")

	v.BindEnv("database.host", "PG_HOST")
	v.BindEnv("database.port", "PG_PORT")
	v.BindEnv("database.user", "PG_USER")
	v.BindEnv("database.password", "PG_PASSWORD")
	v.BindEnv("database.dbname", "PG_DATABASE_NAME")
	v.BindEnv("database.sslmode", "PG_SSL_MODE")

	v.BindEnv("jwt.secret_key", "SECRET_KEY")

	// Значения по умолчанию
	v.SetDefault("server.grpc_port", "50051")
	v.SetDefault("server.host", "localhost")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "54322")
	v.SetDefault("database.user", "auth_db_user")
	v.SetDefault("database.password", "auth_db_password")
	v.SetDefault("database.dbname", "auth_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)

	v.SetDefault("jwt.secret_key", "s12dasd1a3s1d6as5d1a3s1d6as5d")
	v.SetDefault("jwt.access_token_duration", "15m")
	v.SetDefault("jwt.refresh_token_duration", "168h") // 7 дней

	v.SetDefault("rate_limit.requests_per_second", 100)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Error reading config file:", err)
		}
	}
	accessDuration, err := time.ParseDuration(v.GetString("jwt.access_token_duration"))
	if err != nil {
		return nil, fmt.Errorf("invalid access token duration: %v", err)
	}
	refreshDuration, err := time.ParseDuration(v.GetString("jwt.refresh_token_duration"))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token duration: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: v.GetString("server.grpc_port"),
			Host:     v.GetString("server.host"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			DBName:   v.GetString("database.dbname"),
			SSLMode:  v.GetString("database.sslmode"),
			MaxConns: v.GetInt("database.max_conns"),
			MinConns: v.GetInt("database.min_conns"),
		},
		JWT: JWTConfig{
			SecretKey:       v.GetString("jwt.secret_key"),
			AccessTokenTTL:  int(accessDuration.Minutes()),
			RefreshTokenTTL: int(refreshDuration.Minutes()),
		},
		RateLimiter: RateLimiterConfig{
			RequestsPerSecond: v.GetInt("rate_limit.requests_per_second"),
		},
	}
	return cfg, nil
}
