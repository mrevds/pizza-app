package app

import (
	"user-service/internal/handler"
	"user-service/internal/repository/pg"
	"user-service/internal/service"

	"go.uber.org/fx"
)

var Module = fx.Module("app",
	fx.Provide(pg.NewUserRepo),
	fx.Provide(service.NewUserService),
	fx.Provide(handler.NewGRPCHandler),
)
