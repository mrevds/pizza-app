package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mrevds/pizza-app/api-gateway/internal/client"
	"github.com/mrevds/pizza-app/api-gateway/internal/config"
	"github.com/mrevds/pizza-app/api-gateway/internal/logger"
	"github.com/mrevds/pizza-app/api-gateway/internal/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Register() fx.Option {
	return fx.Options(
		fx.Provide(config.Load),
		fx.Provide(func(cfg *config.Config) *zap.Logger {
			return logger.New(cfg.LogLevel).Logger
		}),
		fx.Provide(func(cfg *config.Config) (*client.GRPCClients, error) {
			return client.NewGRPCClients(cfg.UserServiceAddr, cfg.CardServiceAddr)
		}),
		fx.Provide(func(clients *client.GRPCClients, logger *zap.Logger) *http.Server {
			return &http.Server{
				Handler:      router.NewRouter(clients, logger),
				Addr:         fmt.Sprintf(":%d", 8080),
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			}
		}),
		fx.Invoke(runServer),
		fx.Invoke(closeClients),
	)
}

func runServer(lc fx.Lifecycle, server *http.Server, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("server starting", zap.String("addr", server.Addr))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down server")
			return server.Shutdown(ctx)
		},
	})
}

func closeClients(lc fx.Lifecycle, clients *client.GRPCClients, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing gRPC connections")
			return clients.Close()
		},
	})
}
