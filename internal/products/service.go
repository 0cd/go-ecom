package products

import (
	"context"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int64) (repo.Product, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &service{repo: repo}
}

func (s *service) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *service) FindProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductByID(ctx, id)
}
