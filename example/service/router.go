package service

import "github.com/utilslab/iam"

type ShopServiceRouter struct {
	service ShopService
}

func (p ShopServiceRouter) Actions() []iam.Action {
	return []iam.Action{
		{Resources: []iam.Resource{Shop}, Type: iam.Write, Handler: p.service.AddShop, Description: "添加店铺", Codes: AddShopCodes},
		{Resources: []iam.Resource{Shop}, Type: iam.Write, Handler: p.service.AddShop, Description: "添加店铺"},
		{Resources: []iam.Resource{Shop.Optional()}, Type: iam.List, Handler: p.service.AddShop, Description: "添加店铺"},
		{Resources: []iam.Resource{CateGood}, Type: iam.List, Handler: p.service.ListGoods, Description: "获取商品列表"},
	}
}
