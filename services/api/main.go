package main

import (
	"github.com/gin-gonic/gin"

	"xiaodouyin/services/api/handler"
	"xiaodouyin/services/api/jwt"
)

func main() {

	r := gin.Default()
	// 用户服务组
	userGroup := r.Group("douyin/user")
	userGroup.POST("/login/", jwt.AuthMiddleware.LoginHandler)                       // 登陆验证
	userGroup.POST("/register/", handler.CreateUserHandler)                          // 注册
	userGroup.GET("/", jwt.AuthMiddleware.MiddlewareFunc(), handler.UserInfoHandler) // 请求用户数据

	r.Run()
}
