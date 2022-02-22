package service

import (
	"context"
)

type ShopService interface {
	AddShop(ctx context.Context, in AddShopIn) (out AddShopOut, err error)
	ListGoods(ctx context.Context, in AddShopIn) (out AddShopOut, err error)
}

type AddShopIn struct {
}

type AddShopOut struct {
}
