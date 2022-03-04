package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/utilslab/iam"
	"github.com/utilslab/iam/exporter"
	
	service "github.com/utilslab/iam/example/iam/service"
)

func main() {
	
	// gin 跨域配置
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	config.AllowAllOrigins = true
	
	// 自定义 gin 驱动
	engine := gin.Default()
	engine.Use(cors.New(config))
	
	// buck 实例
	api := iam.New()
	api.SetVersion("1.0.0")
	api.SetEngine(engine)
	api.AddRouter(service.NewShopServiceRouter(new(service.Impl)))
	
	// API 导出器配置
	api.SetExporter(":9090", &exporter.Options{
		Project: "Foo",
		Envs: []exporter.Env{
			{
				Name: "本地测试",
				Host: "http://localhost:8080",
			},
		},
	})
	
	api.Run(":8080")
}
