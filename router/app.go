package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ylinyang/im/docs"
	"github.com/ylinyang/im/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	// swagger页面
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 问题
	r.GET("/ping", service.Ping)
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	// 用户
	r.GET("/user-detail", service.GetUserDetail)

	return r
}
