package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/im/models"
	"net/http"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户信息
// @param identity query string false "identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	query := c.Query("identity")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户标识不能为空",
		})
		return
	}

	data := new(models.UserBasic)
	// 查询用户信息，但是排查密码字段
	err := models.DB.Omit("password").Where("identity = ? ", query).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户信息异常" + err.Error(),
		})
		return
	}
	if data.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
