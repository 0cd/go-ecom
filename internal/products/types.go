package products

import "github.com/jackc/pgx/v5/pgtype"

type createProductParams struct {
	Name         pgtype.Text `json:"name"`
	PriceInCents int32       `json:"price_in_cents"`
	Quantity     int32       `json:"quantity"`
}
