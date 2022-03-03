package service

import "github.com/utilslab/iam"

var (
	Shop     = iam.Resource{Name: "shop", Ident: []iam.Field{{Var: "$shopId", Description: "店铺ID"}}, Scope: []iam.Field{{Var: "$color", Description: "颜色"}}, Description: "店铺"}
	Cate     = iam.Resource{Name: "cate", Ident: []iam.Field{{Var: "$cateId", Description: "分类ID"}}, Description: "分类"}
	Good     = iam.Resource{Name: "good", Ident: []iam.Field{{Var: "$goodId", Description: "商品ID"}}, Description: "货物"}
	CateGood = iam.Resource{Name: "CateGood", Ident: []iam.Field{{Var: "$goodId", Description: "分类ID"}}, Description: "分类商品"}
)
