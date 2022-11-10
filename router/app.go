package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/im/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 路由
	r.GET("/ping", service.Ping)
	return r
}
