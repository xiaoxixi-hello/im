package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/im/define"
	"github.com/ylinyang/im/models"
	"net/http"
	"strconv"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @param problemIdentity query string false "problem_identity"
// @param userIdentity query string false "user_identity"
// @param status query string false "status"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-list [get]
func GetSubmitList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	page = (page - 1) * size
	var count int64

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)

	data := make([]models.SubmitBasic, 0)
	// 总共有多少条 + offset 从那一页开始 + limit 每页显示的数据 + 通过find将查询的数据find到data里面
	err := tx.Count(&count).Offset(page).Limit(size).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get submit list err: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  data,
			"count": count,
		},
	})
}
