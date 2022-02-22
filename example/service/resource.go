package service

import "github.com/utilslab/iam"

// 字段权限，可读、可写（一定可读）、隐藏
var (
	Shop     = iam.Resource{Name: "shop", Ident: "shop/$shopId?$color&$size", Description: "店铺"}
	Cate     = iam.Resource{Name: "cate", Ident: "cate/$cateId", Description: "商品"}
	Good     = iam.Resource{Name: "good", Ident: "good/$goodId", Description: "货物"}
	CateGood = iam.Resource{Name: "CateGood", Ident: "$cateId/$goodId", Description: "分类商品"}
)
