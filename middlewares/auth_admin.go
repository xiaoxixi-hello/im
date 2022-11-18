package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/im/define"
	"net/http"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		userClaims, err := define.AnalyseToken(auth)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusExpectationFailed,
				"message": "unauthorized",
			})
			return
		}
		if userClaims.IsAdmin != 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusExpectationFailed,
				"message": "不是管理员",
			})
		}
		c.Next()
	}
}
