package service

import "context"

type Impl struct {
}

func (i Impl) GetShop(ctx context.Context, in AddShopIn) (out AddShopOut, err error) {
	out.ShopId = in.ShopId
	out.CateId = in.CateId
	return
}

func (i Impl) ListGoods(ctx context.Context, in AddShopIn) (out AddShopOut, err error) {
	out.ShopId = in.ShopId
	out.CateId = in.CateId
	return
}
