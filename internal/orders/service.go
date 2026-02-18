package orders

import (
	"context"
	"fmt"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	PlaceOrder(ctx context.Context, order createOrderParams) (repo.Order, error)
	FindOrderByID(ctx context.Context, id int64) (order, error)
}

type service struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &service{
		repo: repo,
		db:   db,
	}
}

func (s *service) PlaceOrder(ctx context.Context, order createOrderParams) (repo.Order, error) {
	if order.UserID == 0 {
		return repo.Order{}, fmt.Errorf("user ID is required")
	}
	if len(order.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, fmt.Errorf("failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	createdOrder, err := qtx.CreateOrder(ctx, order.UserID)
	if err != nil {
		return repo.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range order.Items {
		product, err := qtx.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, fmt.Errorf("product not found: %w", err)
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, fmt.Errorf("product (id: %d) does not have enough stock", item.ProductID)
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      createdOrder.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}

		_, err = qtx.UpdateProduct(ctx, repo.UpdateProductParams{
			ID: product.ID,
			Quantity: pgtype.Int4{
				Int32: product.Quantity - item.Quantity,
				Valid: true,
			},
		})
		if err != nil {
			return repo.Order{}, fmt.Errorf("failed to update product quantity: %w", err)
		}
	}

	tx.Commit(ctx)

	return createdOrder, nil
}

func (s *service) FindOrderByID(ctx context.Context, id int64) (order, error) {
	foundOrders, err := s.repo.FindOrderByID(ctx, id)
	if err != nil {
		return order{}, fmt.Errorf("failed to find order: %w", err)
	}

	if len(foundOrders) == 0 {
		return order{}, fmt.Errorf("order not found")
	}

	items := make([]orderItem, 0, len(foundOrders))
	var totalPrice int32
	for _, row := range foundOrders {
		items = append(items, orderItem{
			ProductID:    row.ProductID,
			Quantity:     row.Quantity,
			PriceInCents: row.PriceInCents,
		})

		totalPrice += row.Quantity * row.PriceInCents
	}

	o := order{
		ID:         foundOrders[0].ID,
		UserID:     foundOrders[0].UserID,
		Items:      items,
		TotalPrice: totalPrice,
		CreatedAt:  foundOrders[0].CreatedAt,
	}

	return o, nil
}
