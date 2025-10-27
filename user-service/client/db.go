package client

import (
	"user-service/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
)

type DB struct {
	Pool *pgxpool.Pool
}

// NewDB создает подключение к БД с lifecycle hooks
func NewDB(lc fx.Lifecycle, cfg *config.Config) (*DB, error) {
	// Парсим конфигурацию
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Настройка пула
	poolConfig.MaxConns = cfg.Database.MaxConns
	poolConfig.MinConns = cfg.Database.MinConns
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	var pool *pgxpool.Pool
	db := &DB{}

	// Lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			if err = pool.Ping(ctx); err != nil {
				pool.Close()
				return fmt.Errorf("failed to ping database: %w", err)
			}

			db.Pool = pool
			fmt.Println("✅ Database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if db.Pool != nil {
				db.Pool.Close()
				fmt.Println("✅ Database connection closed")
			}
			return nil
		},
	})

	return db, nil
}