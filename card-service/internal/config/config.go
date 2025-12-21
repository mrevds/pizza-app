package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   CardServerConfig
	Database CardDataBaseConfig
}

type CardServerConfig struct {
	GRPCPort string
	Host     string
}

type CardDataBaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

func CardConfigLoad() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/card-service/")

	v.AutomaticEnv()
	v.BindEnv("server.grpc_port", "GRPC_PORT")
	v.BindEnv("server.host", "SERVER_HOST")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5433")
	v.SetDefault("database.user", "card_db_user")
	v.SetDefault("database.password", "card_db_password")
	v.SetDefault("database.dbname", "card_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Error reading config file: %w", err)
		}
	}
	cfg := &Config{
		Server: CardServerConfig{
			GRPCPort: v.GetString("server.grpc_port"),
			Host:     v.GetString("server.host"),
		},
		Database: CardDataBaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			DBName:   v.GetString("database.dbname"),
			SSLMode:  v.GetString("database.sslmode"),
			MaxConns: v.GetInt("database.max_conns"),
			MinConns: v.GetInt("database.min_conns"),
		},
	}
	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
	c.Database.Host,
	c.Database.Port,
	c.Database.User,
	c.Database.Password,
	c.Database.DBName,
	c.Database.SSLMode,
	)
}
