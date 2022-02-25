package service

import (
	"context"
)

type ShopService interface {
	GetShop(ctx context.Context, in AddShopIn) (out AddShopOut, err error)
	ListGoods(ctx context.Context, in AddShopIn) (out AddShopOut, err error)
}

type AddShopIn struct {
	ShopId int64
	CateId int64
}

type AddShopOut struct {
}
