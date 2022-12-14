package main

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/fkssl/apis"
	"github.com/suisrc/fkssl/serve"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 配置路由
	router.GET("/ping", serve.Ping)
	router.GET("/healthz", serve.Healthz)
	// 注册业务路由
	router.POST("/api/ssl/v1/ca/init", apis.CreateCaCmdApi)
	router.GET("/api/ssl/v1/ca", apis.QuaryCaQryApi)
	router.GET("/api/ssl/v1/ca/txt", apis.QuaryCaQryTxtApi)
	router.GET("/api/ssl/v1/cert", apis.QurayCertCmdApi)
	// 启动服务
	serve.Run(router)
}
