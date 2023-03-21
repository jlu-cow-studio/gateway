package main

import "github.com/jlu-cow-studio/gateway/handler"

func RegisterHandlers() {
	server.GET("/ping", handler.Ping)

	server.POST("/user/login", handler.UserLogin)
	server.POST("/user/register", handler.UserRegister)
	server.POST("/user/auth", handler.UserAuth)
	server.GET("/user/info", handler.UserInfo)
}
