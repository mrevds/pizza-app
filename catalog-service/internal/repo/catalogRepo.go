package repo

import "catalog-service/internal/database"

type catalogRepo struct {
	db *database.DB
}

func NewCatalogRepo(db *database.DB) CatalogRepo {
	return &catalogRepo{db: db}
}
