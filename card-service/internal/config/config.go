package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	DataBase DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Host string
	Port int
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

//type RateLimiterConfig struct {
//	RequestsPerMinute int
//}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/card-service/")
	v.AutomaticEnv()

	v.BindEnv("server.grpc_host", "GRPC_HOST")
	v.BindEnv("server.grpc_port", "GRPC_PORT")

	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.dbname", "DB_NAME")
	v.BindEnv("database.sslmode", "DB_SSLMODE")
	v.BindEnv("database.max_conns", "DB_MAX_CONNS")
	v.BindEnv("database.min_conns", "DB_MIN_CONNS")
	v.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	v.BindEnv("jwt.access_token_ttl", "JWT_ACCESS_TOKEN_TTL")
	v.BindEnv("jwt.refresh_token_ttl", "JWT_REFRESH_TOKEN_TTL")

	v.SetDefault("server.grpc_port", "50052")
	v.SetDefault("server.host", "localhost")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "54333")
	v.SetDefault("database.user", "user_db_user")
	v.SetDefault("database.password", "user_db_password")
	v.SetDefault("database.dbname", "user_db")
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
			Port: v.GetInt("server.grpc_port"),
			Host: v.GetString("server.host"),
		},
		DataBase: DatabaseConfig{
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
	}
	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DataBase.Host,
		c.DataBase.Port,
		c.DataBase.User,
		c.DataBase.Password,
		c.DataBase.DBName,
		c.DataBase.SSLMode,
	)
}
