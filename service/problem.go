package service

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/ylinyang/im/define"
	"github.com/ylinyang/im/models"
	"gorm.io/gorm"
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
// @param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	// 获取接口传入参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")
	// page 1 == > db 0
	page = (page - 1) * size
	var count int64

	list := make([]models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @param identity query string false "identity"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	problemDetail := c.Query("identity")
	if problemDetail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题详情不能为空",
		})
		return
	}
	m := new(models.ProblemBasic)
	err := models.DB.Where("identity = ? ", problemDetail).
		// 需要连表查询
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取问题详情异常" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": m,
	})
	return
}

// ProblemCreate
// @Tags 私有方法
// @Summary 管理员问题创建
// @param token header string true "token"
// @param title formData string true "title"
// @param content formData string true "content"
// @param max_runtime formData int false "max_runtime"
// @param max_mem formData int false "max_mem"
// @param category_ids formData array true "test_cases"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /problem-create [post]
func ProblemCreate(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := uuid.NewV4().String()
	data := models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}
	// 处理分类

	// 处理测试用例

	if models.DB.Create(&data).Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户问题表创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})
}
