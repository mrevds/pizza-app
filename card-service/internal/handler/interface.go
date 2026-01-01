package handler

import (
	"context"
	pbCard "card-service/pkg/card_v1/card-service_v1"
)

type GRPCHandler struct {
	pbCard.CardServiceServer
}

type HealthCheck interface {
	Ping(ctx context.Context) error
}