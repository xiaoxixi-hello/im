package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ylinyang/im/define"
	"github.com/ylinyang/im/models"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	// 获取接口传入参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	keyword := c.Query("keyword")
	// page 1 == > db 0
	page = (page - 1) * size
	var count int64

	list := make([]models.Problem, 0)
	tx := models.GetProblemList(keyword)
	err := tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("get problem list Error: ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"list": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
