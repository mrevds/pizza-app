package handler

import (
	pbCatalog "catalog-service/pkg/catalog_v1/catalog-service_v1"
	"context"
)

type GRPCHandler struct {
	pbCatalog.CatalogServiceServer
}

type HealthCheck interface {
	Ping(ctx context.Context) error
}
