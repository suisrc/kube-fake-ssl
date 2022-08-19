package main

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/fkssl/serve"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 配置路由
	router.GET("/ping", serve.Ping)
	router.GET("/healthz", serve.Healthz)
	// 启动服务
	serve.Run(router)
}
