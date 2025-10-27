package handler

import (
	"context"
	userGRPC "user-service/pkg/user-service_v1"
)

type GRPCHandler struct {
	userGRPC.UserServiceServer
}

type HealthCheck interface {
	Ping(ctx context.Context) error
}
