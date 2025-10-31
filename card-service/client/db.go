package client

import (
	"context"
	"fmt"
	"github.com/mrevds/pizza-app/card-service/internal/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(lc fx.Lifecycle, cfg *config.Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}
	poolConfig.MaxConns = int32(cfg.DataBase.MaxConns)
	poolConfig.MinConns = int32(cfg.DataBase.MinConns)
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	var pool *pgxpool.Pool
	db := &DB{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
			if err != nil {
				return fmt.Errorf("Error", err)
			}
			if err = pool.Ping(ctx); err != nil {
				pool.Close()
				return fmt.Errorf("failed to ping database", err)
			}
			db.Pool = pool
			fmt.Println("Database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if db.Pool != nil {
				db.Pool.Close()
				fmt.Println("Database connection closed")
			}
			return nil
		},
	})
	return db, err
}
