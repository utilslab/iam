package service

import "github.com/utilslab/iam"

var AddShopCodes = []iam.Code{
	{Status: 400, Code: "GoodDuplicate", Message: "商品已经存在"},
}
