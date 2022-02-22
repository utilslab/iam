package service

import (
	"github.com/utilslab/iam"
)

func NewFooRouter(service FooService) *FooRouter {
	return &FooRouter{service: service}
}

type FooRouter struct {
	service FooService
}

func (f FooRouter) Routes() []iam.Route {
	return []iam.Route{
		{Method: iam.Get, Handler: f.service.Ping, Description: "测试", Extra: "customer(color,tag):customerId,cluster(bool,mini):clusterId"},
		{Method: iam.Get, Handler: f.service.GetHtml},
		{Method: iam.Get, Handler: f.service.GetText},
		{Method: iam.Get, Handler: f.service.GetInt},
		{Method: iam.Get, Handler: f.service.GetInt32},
		{Method: iam.Get, Handler: f.service.QueryPost},
		{Method: iam.Get, Handler: f.service.GetDecimal},
		{Method: iam.Get, Handler: f.service.GetBool},
		{Handler: f.service.AddPost, Extra: "customerId,clusterId"},
		{Handler: f.service.TestGetArray},
		{Handler: f.service.TestPostArray},
		{Handler: f.service.Ping2},
		{Handler: f.service.PostShop},
	}
}
