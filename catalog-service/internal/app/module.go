package app

import (
	"catalog-service/internal/config"
	"catalog-service/internal/handler"
	"catalog-service/internal/middleware"
	"catalog-service/internal/repo"
	"catalog-service/internal/service"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func newCatalogGRPCServer(cfg *config.Config) *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTInterceptor(cfg.JWT.Secret)),
	)
}

var CatalogModule = fx.Module("app",
	fx.Provide(repo.NewCatalogRepo),
	fx.Provide(service.NewCatalogService),
	fx.Provide(handler.NewCatalogGRPCHandler),
	fx.Provide(newCatalogGRPCServer),
)
