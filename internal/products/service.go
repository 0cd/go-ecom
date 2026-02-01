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
	UpdateProduct(ctx context.Context, updates repo.UpdateProductParams) (repo.Product, error)
	ReplaceProduct(ctx context.Context, newProduct repo.UpdateProductParams) (repo.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
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
		return repo.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	return createdProduct, nil
}

func (s *service) UpdateProduct(ctx context.Context, updates repo.UpdateProductParams) (repo.Product, error) {
	if err := s.validateProductUpdates(updates, false); err != nil {
		return repo.Product{}, err
	}

	updatedProduct, err := s.repo.UpdateProduct(ctx, updates)
	if err != nil {
		return repo.Product{}, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

func (s *service) ReplaceProduct(ctx context.Context, newProduct repo.UpdateProductParams) (repo.Product, error) {
	if err := s.validateProductUpdates(newProduct, true); err != nil {
		return repo.Product{}, err
	}

	updatedProduct, err := s.repo.UpdateProduct(ctx, newProduct)
	if err != nil {
		return repo.Product{}, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

func (s *service) DeleteProduct(ctx context.Context, id int64) error {
	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *service) validateProductUpdates(updates repo.UpdateProductParams, requireAll bool) error {
	_, err := s.repo.FindProductByID(context.Background(), updates.ID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}

	if requireAll {
		if !updates.Name.Valid || updates.Name.String == "" {
			return fmt.Errorf("name is required")
		}
		if !updates.PriceInCents.Valid {
			return fmt.Errorf("price is required")
		}
		if !updates.Quantity.Valid {
			return fmt.Errorf("quantity is required")
		}
	}

	if updates.Name.Valid && updates.Name.String == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if updates.PriceInCents.Valid && updates.PriceInCents.Int32 < 0 {
		return fmt.Errorf("price must be equal or greater than 0")
	}

	if updates.Quantity.Valid && updates.Quantity.Int32 < 0 {
		return fmt.Errorf("quantity must be equal or greater than 0")
	}

	return nil
}
