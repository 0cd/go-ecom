package orders

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type orderItem struct {
	ProductID    int64 `json:"product_id"`
	Quantity     int32 `json:"quantity"`
	PriceInCents int32 `json:"price_in_cents"`
}

type createOrderParams struct {
	UserID int64       `json:"user_id"`
	Items  []orderItem `json:"items"`
}

type order struct {
	ID         int64              `json:"id"`
	UserID     int64              `json:"user_id"`
	Items      []orderItem        `json:"items"`
	TotalPrice int32              `json:"total_price"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
}
