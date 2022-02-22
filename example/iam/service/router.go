package service

import "github.com/utilslab/iam"

var _ iam.Service = &ShopServiceRouter{}

type ShopServiceRouter struct {
	service ShopService
}

func (p ShopServiceRouter) Actions() []iam.ActionGroup {
	return []iam.ActionGroup{
		{
			Name: "商品",
			Actions: []iam.Action{
				{Resources: []iam.Resource{Shop}, Type: iam.Write, Handle: p.service.AddShop, Description: "添加店铺", Codes: AddShopCodes},
				{Resources: []iam.Resource{Shop}, Type: iam.Write, Handle: p.service.AddShop, Description: "添加店铺"},
				{Resources: []iam.Resource{Shop.Optional()}, Type: iam.List, Handle: p.service.AddShop, Description: "添加店铺"},
			},
		},
		{
			Name: "分类商品",
			Actions: []iam.Action{
				{Resources: []iam.Resource{CateGood}, Type: iam.List, Handle: p.service.ListGoods, Description: "获取商品列表"},
			},
		},
	}
}
