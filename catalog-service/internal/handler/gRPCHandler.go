package handler

import (
	"catalog-service/internal/service"
	pbCatalog "catalog-service/pkg/catalog_v1/catalog-service_v1"
)

type grpcHandler struct {
	catalogServ service.CatalogService
	pbCatalog.UnimplementedCatalogServiceServer
}

func NewCatalogGRPCHandler(catalogServ service.CatalogService) pbCatalog.CatalogServiceServer {
	return &grpcHandler{
		catalogServ: catalogServ,
	}
}
