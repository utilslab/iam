package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/utilslab/iam"
	"github.com/utilslab/iam/example/test/service"
	"github.com/utilslab/iam/exporter"
)

func main() {
	
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	config.AllowAllOrigins = true
	
	engine := gin.Default()
	//engine.Use(cors.New(config), auth.Guard())
	
	api := iam.New()
	api.SetVersion("1.0.0")
	api.AddRouter(service.NewUserRouter(new(service.UserImplService)))
	api.SetExporter(":8090", &exporter.Options{
		Project: "测试项目",
		Envs: []exporter.Env{
			{
				Name: "本地测试",
				Host: "http://localhost:8088",
			},
			{
				Name: "另外一个测试环境",
				Host: "http://localhost:8090",
			},
		},
	})
	api.SetEngine(engine)
	api.Run(":8088")
}
