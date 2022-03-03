package service

import (
	"github.com/gin-gonic/gin"
	"github.com/utilslab/iam"
)

//var _ iam.Service = &ShopServiceRouter{}

func NewShopServiceRouter(service ShopService) *ShopServiceRouter {
	return &ShopServiceRouter{service: service}
}

type ShopServiceRouter struct {
	service ShopService
}

func (p ShopServiceRouter) Routes() []*iam.Route {
	return []*iam.Route{
		{
			Prefix:      "/api",
			Middlewares: []gin.HandlerFunc{},
			Groups: []*iam.Group{
				{
					Name:   "商品",
					Prefix: "/good",
					Actions: []*iam.Action{
						// In -> $ShopId, $CateId
						{Resources: []iam.Resource{Shop, Cate}, Type: iam.Read, Handler: p.service.GetShop, Description: "获取店铺", Codes: AddShopCodes},
					},
				},
				{
					Name:   "分类商品",
					Prefix: "/cateGood",
					Actions: []*iam.Action{
						{Resources: []iam.Resource{CateGood}, Type: iam.List, Handler: p.service.ListGoods, Description: "获取商品列表"},
					},
				},
			},
		},
	}
}
