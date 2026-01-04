package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	JWT           JWTConfig
	KafkaProducer KafkaConfig
}

type JWTConfig struct {
	Secret string
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

type KafkaConfig struct {
	Brokers []string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/catalog-service/")

	v.AutomaticEnv()
	v.BindEnv("server.grpc_port", "GRPC_PORT")
	v.BindEnv("server.host", "SERVER_HOST")
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.dbname", "DB_NAME")
	v.BindEnv("jwt.secret", "JWT_SECRET")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5434")
	v.SetDefault("database.user", "catalog_db_user")
	v.SetDefault("database.password", "catalog_db_password")
	v.SetDefault("database.dbname", "catalog_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)
	v.SetDefault("kafka.brokers", "localhost:9092")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
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
			Secret: v.GetString("jwt.secret"),
		},
		KafkaProducer: KafkaConfig{
			Brokers: v.GetStringSlice("kafka.brokers"),
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
