package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ylinyang/im/docs"
	"github.com/ylinyang/im/middlewares"
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
	r.POST("/register", service.Register)
	r.POST("/login", service.Login)
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/send-email", service.SendMail)
	r.GET("/user-rank-list", service.UserRankList)

	// 提交记录
	r.GET("/submit-list", service.GetSubmitList)

	// 下面的路由只有管理员才有权限 需要先进行验证
	r.POST("/problem-create", middlewares.AuthAdminCheck(), service.ProblemCreate)

	return r
}
