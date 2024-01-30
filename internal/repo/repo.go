package repo

import (
	"context"

	"github.com/KapDmitry/WB_L0/internal/order"
)

type Repo interface {
	GetAll(ctx context.Context) ([]order.Order, error)
	Add(ctx context.Context, ord order.Order) error
}
