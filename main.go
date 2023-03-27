package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/jlu-cow-studio/gateway/handler"
	"github.com/jlu-cow-studio/gateway/middleware"
)

const ServerAddress = "0.0.0.0:8080"

var server *gin.Engine

func main() {
	discovery.Init()
	redis.Init()

	server = gin.New()

	server.Use(middleware.RequestCheck()) //请求校验filter
	server.Use(middleware.GetUserInfo())  //用户信息filter

	RegisterHandlers()

	server.Run(ServerAddress)
}

func RegisterHandlers() {
	server.GET("/ping", handler.Ping)

	server.POST("/user/login", handler.UserLogin)
	server.POST("/user/register", handler.UserRegister)
	server.POST("/user/auth", handler.UserAuth)
	server.GET("/user/info", handler.UserInfo)

	server.POST("/item/add", handler.ItemAdd)
	server.DELETE("/item/delete", handler.ItemDelete)
	server.POST("/item/update", handler.ItemUpdate)

	server.GET("/feed/v1", handler.Feedv1)
	server.PUT("/feed/v1", handler.Feedv1Reset)
}
