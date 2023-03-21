package main

import "github.com/jlu-cow-studio/gateway/handler"

func RegisterHandlers() {
	server.GET("/ping", handler.Ping)

	server.POST("/user/login", handler.UserLogin)
}
