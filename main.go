package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/jlu-cow-studio/gateway/middleware"
)

const ServerAddress = "0.0.0.0:8080"

func main() {
	discovery.Init()
	redis.Init()

	r := gin.New()

	r.Use(middleware.RequestCheck()) //请求校验filter
	r.Use(middleware.GetUserInfo())  //用户信息filter

	r.GET("/ping", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "pong")
	})

	r.Run(ServerAddress)
}
