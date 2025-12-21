package main

import (
	"card-service/database"
	"card-service/internal/config"
	pbCard "card-service/pkg/card_v1/card-service_v1"
	"context"
	"fmt"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fx.New(
		fx.Provide(
			config.CardConfigLoad,
			database.CardDbInit,
		),
		//app.Module, надо написать эту дрочь DI
		fx.Invoke(registergRPCServer),
	).Run()
}

func registergRPCServer(
	lc fx.Lifecycle,
	grpcServer *grpc.Server,
	handler pbCard.CardServiceServer,
	cfg *config.Config) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				lis, err := net.Listen("tcp", fmt.Sprintf(":%s",cfg.Server.GRPCPort))
				if err != nil {
					return err
				}
				pbCard.RegisterCardServiceServer(grpcServer, handler)
				reflection.Register(grpcServer)
				go func() {
					log.Printf("gRPC server listenint at: %s", cfg.Server.GRPCPort)
					if err := grpcServer.Serve(lis); err != nil {
						log.Fatalf("failed to serve: %v", err)
					}
				}()
				return nil
			},
		})
	}
