package app

import (
	"card-service/internal/config"
	"card-service/internal/handler"
	"card-service/internal/repo/pg"
	"card-service/internal/service"
	"card-service/middleware"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func newCardGRPCServer(cfg *config.Config) *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTInterceptor(cfg.JWT.Secret)),
	)
}

var Cardmodule = fx.Module("app",
	fx.Provide(pg.NewCardRepo),
	fx.Provide(service.NewCardService),
	fx.Provide(handler.NewCardGRPCHandler),
	fx.Provide(newCardGRPCServer),
)
