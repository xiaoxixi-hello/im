# im
设计思路
1. 设计数据库  创建models文件夹，存储数据表 一张表对应一个go文件， 初始化数据库连接
2. 路由  配置路由统一存放位置
3. swagger使用 -- 安装官网配置好swagger 进入http://127.0.0.1:8080/swagger/index.html
   https://github.com/swaggo/gin-swagger
4. 完成问题接口
5. 用户登录接口 生成token 使用jwt
   https://pkg.go.dev/github.com/dgrijalva/jwt-go#example-Parse--Hmac
6. 使用邮箱发送验证码功能
   https://github.com/jordan-wright/email
7. 用户注册
   UUID 用户唯一标识  https://github.com/satori/go.uuid
   redis 缓存  验证码过期 https://github.com/go-redis/redis
   