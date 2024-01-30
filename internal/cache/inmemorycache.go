package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/KapDmitry/WB_L0/internal/order"
	"github.com/KapDmitry/WB_L0/internal/repo"
)

type InMemoryCash struct {
	Mu   *sync.Mutex
	Cash map[string]order.Order
}

func NewInMemoryCash() *InMemoryCash {
	return &InMemoryCash{
		Mu:   &sync.Mutex{},
		Cash: make(map[string]order.Order),
	}
}

func (i *InMemoryCash) Get(orderID string) (order.Order, error) {
	i.Mu.Lock()
	defer i.Mu.Unlock()
	if val, ok := i.Cash[orderID]; ok {
		return val, nil
	}
	return order.Order{}, fmt.Errorf("no such order")
}

func (i *InMemoryCash) Update(order order.Order) error {
	i.Mu.Lock()
	defer i.Mu.Unlock()
	if _, ok := i.Cash[order.OrderUID]; ok {
		return fmt.Errorf("order already exists")
	}
	i.Cash[order.OrderUID] = order
	return nil
}

func (i *InMemoryCash) Recover(ctx context.Context, rep repo.Repo) error {
	orders, err := rep.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, order := range orders {
		err = i.Update(order)
		if err != nil {
			return err
		}
	}
	return nil
}
