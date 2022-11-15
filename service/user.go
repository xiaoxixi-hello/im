package service

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/ylinyang/im/define"
	"github.com/ylinyang/im/models"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

// Login
// @Tags 公共方法
// @Summary 用户登录
// @param username formData string false "username"
// @param password formData string false "password"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户密码不能为空",
		})
		return
	}

	// 将密码通过md5加密 %x 小写
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	log.Println(username)
	log.Println(password)

	data := new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或者密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询用户未知异常," + err.Error(),
		})
		return
	}
	token, err := define.GenerateToken(data.Identity, username, data.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成token异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})

}

// SendMail
// @Tags 公共方法
// @Summary 发送验证码
// @param email formData string true "email"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /send-email [post]
func SendMail(c *gin.Context) {
	mail := c.PostForm("email")
	// 生成随机验证码
	code := ""
	// 生成随机数种子
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10)) // 0-10随便选一个
	}

	// 发送邮件
	e := email.NewEmail()
	e.From = "master <2025907338@qq.com>"
	e.To = []string{mail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	err := e.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "2025907338@qq.com", define.EmailPassWord, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发送邮件异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发送邮件成功",
	})
	// 将发送的验证码存入redis
	define.RedisDB.Set(context.TODO(), mail, code, time.Second*60)
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @param username formData string true "username"
// @param password formData string true "password"
// @param code formData string true "code"
// @param mail formData string true "mail"
// @param phone formData string false "phone"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	code := c.PostForm("code")
	mail := c.PostForm("mail")
	phone := c.PostForm("phone")
	if username == "" || password == "" || code == "" || mail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	// 验证码 存入redis 10s过期时间
	redisCode, err := define.RedisDB.Get(context.Background(), mail).Result()
	if err != nil {
		log.Println("从redis获取验证码失败")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取验证码失败,请重试",
		})
		return
	}
	if code != redisCode {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "验证码不正确",
		})
	}
	// 验证码正确 判断mail是否重复
	var eCount int64
	if models.DB.Where("mail = ? ", mail).Model(new(models.UserBasic)).Count(&eCount).Error != nil {
		log.Println("查询改用户是否注册失败")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户注册失败,请重试",
		})
		return
	}
	if eCount > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该用户已经注册",
		})
		return
	}

	//用户信息插入  生成token
	data := &models.UserBasic{
		Identity: uuid.NewV4().String(),
		Name:     username,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(password))),
		Phone:    phone,
		Mail:     mail,
		IsAdmin:  0,
	}
	if models.DB.Create(data).Error != nil {
		log.Println("用户信息插入数据库失败")
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户注册失败,请重试",
		})
		return
	}
	token, _ := define.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

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

// UserRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":200,"data":""}"
// @Router /user-rank-list [get]
func UserRankList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page = (page - 1) * size

	var userRankListCount int64
	list := make([]models.UserBasic, 0)
	// order 按照字段进行排序
	if models.DB.Model(new(models.UserBasic)).Count(&userRankListCount).Order("pass_num DESC, submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取排行榜失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": userRankListCount,
		},
	})
}
