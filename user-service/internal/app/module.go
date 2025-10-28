package app

import (
	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repository/pg"
	"user-service/internal/service"
	"user-service/internal/utils"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func newGRPCServer(authInterceptor *middleware.AuthInterceptor) *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
}

var Module = fx.Module("app",
	fx.Provide(pg.NewUserRepo),
	fx.Provide(service.NewUserService),
	fx.Provide(handler.NewGRPCHandler),
	fx.Provide(newGRPCServer),
	fx.Provide(utils.NewJWTManager),
	fx.Provide(middleware.NewAuthInterceptor),
)
