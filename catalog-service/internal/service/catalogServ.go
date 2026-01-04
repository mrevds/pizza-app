package service

import "catalog-service/internal/repo"

type catalogService struct {
	catalogRepo repo.CatalogRepo
}

func NewCatalogService(catalogRepo repo.CatalogRepo) CatalogService {
	return &catalogService{
		catalogRepo: catalogRepo,
	}
}
