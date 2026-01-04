package database

import (
	"catalog-service/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}

func CatalogDbInit(lc fx.Lifecycle, cfg *config.Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}
	poolConfig.MaxConns = int32(cfg.Database.MaxConns)
	poolConfig.MinConns = int32(cfg.Database.MinConns)
	poolConfig.MaxConnIdleTime = time.Hour

	var pool *pgxpool.Pool
	db := &DB{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
			if err != nil {
				return fmt.Errorf("pgxpool.ConnectConfig: %w", err)
			}
			if err = pool.Ping(ctx); err != nil {
				return fmt.Errorf("ping: %w", err)
			}
			db.Pool = pool
			fmt.Println("Connected to database")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if db.Pool != nil {
				db.Pool.Close()
				fmt.Println("Disconnected from database")
			}
			return nil
		},
	})
	return db, nil
}
