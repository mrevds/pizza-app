package main

import (
	"fmt"
	"log"
	"net"
	"user-service/client"
	"user-service/internal/app"
	"user-service/internal/config"
	"user-service/internal/handler"

	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "user-service/pkg/user-service_v1"
)

func main() {
	fx.New(
		fx.Provide(
			config.Load,
			client.NewDB,
			handler.NewGRPCHandler,
		),
		app.Module,
		fx.Invoke(registerGRPServer),
	).Run()
}

func registerGRPServer(
	lc fx.Lifecycle,
	grpcServer *grpc.Server,
	handler pb.UserServiceServer,
	cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
			if err != nil {
				return err
			}
			pb.RegisterUserServiceServer(grpcServer, handler)
			reflection.Register(grpcServer)
			go func() {
				log.Printf("GRPC server listening at %s", cfg.Server.GRPCPort)
				if err := grpcServer.Serve(lis); err != nil {
					log.Fatalf("failed to serve: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping GRPC server...")
			grpcServer.GracefulStop()
			log.Printf("GRPC server stopped.")
			return nil
		},
	})
}
