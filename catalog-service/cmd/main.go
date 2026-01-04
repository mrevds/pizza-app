package main

import (
	"catalog-service/internal/app"
	"catalog-service/internal/config"
	"catalog-service/internal/database"
	pbCatalog "catalog-service/pkg/catalog_v1/catalog-service_v1"
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	fx.New(
		fx.Provide(
			config.Load,
			database.CatalogDbInit,
		),
		app.CatalogModule,
		fx.Invoke(registerGRPCServer),
	).Run()
}

func registerGRPCServer(
	lc fx.Lifecycle,
	grpcServer *grpc.Server,
	handler pbCatalog.CatalogServiceServer,
	cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
			if err != nil {
				return err
			}
			pbCatalog.RegisterCatalogServiceServer(grpcServer, handler)
			reflection.Register(grpcServer)
			go func() {
				log.Printf("gRPC server listening on port %s", cfg.Server.GRPCPort)
				if err := grpcServer.Serve(lis); err != nil {
					log.Fatalf("gRPC server failed to serve: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("gRPC server stopping")
			grpcServer.GracefulStop()
			log.Printf("gRPC server stopped")
			return nil
		},
	})
}
