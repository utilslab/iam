package service

import (
	"github.com/utilslab/iam"
)

func NewUserRouter(service UserService) *UserRouter {
	return &UserRouter{service: service}
}

type UserRouter struct {
	service UserService
}

func (p UserRouter) Routes() []iam.Route {
	return []iam.Route{
		{
			Prefix: "/api",
			Children: []iam.Route{
				{
					Prefix: "",
					Children: []iam.Route{
						{Method: iam.Get, Path: "/ping", Handler: p.service.Ping, Description: "测试项目是否连通"},
						{Method: iam.Get, Handler: p.service.Wait},
						{Method: iam.Get, Path: "/inner-error", Handler: p.service.InnerError},
						{Method: iam.Get, Handler: p.service.ValidateError},
						{Method: iam.Get, Handler: p.service.ForbiddenError},
					},
				},
				{
					Prefix: "/guard",
					Children: []iam.Route{
						{Handler: p.service.TestPost},
						{Handler: p.service.TestPostArray},
						{Method: iam.Get, Handler: p.service.TestGet},
						{Method: iam.Get, Handler: p.service.TestGetArray},
						{Method: iam.Put, Handler: p.service.TestPut},
						{Method: iam.Put, Handler: p.service.TestPutArray},
						{Method: iam.Delete, Handler: p.service.TestDelete},
						{Method: iam.Delete, Handler: p.service.TestDeleteArray},
						{Method: iam.Get, Handler: p.service.TestDecimal},
						{Method: iam.Post, Handler: p.service.TestNestedInput},
					},
				},
			},
		},
	}
}
