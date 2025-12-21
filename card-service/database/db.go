package database

import (
	"card-service/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
)

type DB struct {
	Pool *pgxpool.Pool
}

func CardDbInit(lc fx.Lifecycle, cfg *config.Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error in carddb init: %w", err)
	}
	poolConfig.MaxConns = int32(cfg.Database.MaxConns)
	poolConfig.MinConns = int32(cfg.Database.MinConns)
	poolConfig.MaxConnIdleTime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	var pool *pgxpool.Pool
	db := &DB{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
			if err != nil {
				return fmt.Errorf("error conn in lc db: %w", err)
			}
			if err = pool.Ping(ctx); err != nil {
				pool.Close()
				return fmt.Errorf("failed to ping db: %w", err)
			}
			db.Pool = pool
			fmt.Println("db conn success")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if db.Pool != nil {
				db.Pool.Close()
				fmt.Println("db conn closed")
			}
			return nil
		},
	})
	return db, nil
}
