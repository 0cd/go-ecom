package products

import (
	"context"
	"fmt"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, product createProductParams) (repo.Product, error)
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

func (s *service) CreateProduct(ctx context.Context, product createProductParams) (repo.Product, error) {
	if product.Name.String == "" {
		return repo.Product{}, fmt.Errorf("name is required")
	}
	if product.PriceInCents < 0 {
		return repo.Product{}, fmt.Errorf("price must be equal or greater than 0")
	}
	if product.Quantity < 0 {
		return repo.Product{}, fmt.Errorf("quantity must be equal or greater than 0")
	}

	createdProduct, err := s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         product.Name.String,
		PriceInCents: product.PriceInCents,
		Quantity:     product.Quantity,
	})
	if err != nil {
		return repo.Product{}, fmt.Errorf("failed to create order: %w", err)
	}

	return createdProduct, nil
}
