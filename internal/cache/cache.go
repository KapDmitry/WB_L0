package cache

import (
	"context"

	"github.com/KapDmitry/WB_L0/internal/order"
	"github.com/KapDmitry/WB_L0/internal/repo"
)

type Cache interface {
	Get(orderID string) (order.Order, error)
	Update(order order.Order) error
	Recover(ctx context.Context, rep repo.Repo) error
}
