package orderuc

import (
	"context"
	"fmt"

	"mimokocke/internal/domain/order"
)

type CreateOrderUseCase struct {
	orderSrv *order.Service
}

func NewCreateOrderUseCase(orderSrv *order.Service) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderSrv: orderSrv,
	}
}

func (c *CreateOrderUseCase) Execute(ctx context.Context) {
	fmt.Println("create order")

	// Transaction
}
