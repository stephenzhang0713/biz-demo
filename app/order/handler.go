package main

import (
	"context"
	"github.com/baiyutang/gomall/app/order/biz/service"
	order "github.com/baiyutang/gomall/app/order/kitex_gen/order"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// PlaceOrder implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderRequest) (resp *order.PlaceOrderResponse, err error) {
	resp, err = service.NewPlaceOrderService(ctx).Run(req)

	return resp, err
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderRequest) (resp *order.ListOrderResponse, err error) {
	resp, err = service.NewListOrderService(ctx).Run(req)

	return resp, err
}
