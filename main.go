package main

import (
	"os"

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

	server.Use(middleware.CROS())         //跨域请求
	server.Use(middleware.RequestCheck()) //请求校验filter
	server.Use(middleware.GetUserInfo())  //用户信息filter

	if err := os.MkdirAll("/opt/img", 0755); err != nil {
		panic(err)
	}

	RegisterHandlers()

	server.Run(ServerAddress)
}

func RegisterHandlers() {
	server.GET("/ping", handler.Ping)

	server.POST("/user/login", handler.UserLogin)
	server.POST("/user/register", handler.UserRegister)
	server.POST("/user/auth", handler.UserAuth)
	server.POST("/user/info", handler.UserInfo)

	server.POST("/user/follow", handler.UserFollow)
	server.POST("/user/following", handler.UserFollowing)
	server.POST("/user/follower", handler.UserFollower)
	server.POST("/user/follow_count", handler.UserFollowCount)

	server.POST("/item/add", handler.ItemAdd)
	server.DELETE("/item/delete", handler.ItemDelete)
	server.POST("/item/update", handler.ItemUpdate)
	server.POST("/item/add_favorite", handler.AddFavorite)

	server.POST("/feed/v1", handler.Feedv1)
	server.PUT("/feed/v1", handler.Feedv1Reset)

	server.POST("/trade/recharge", handler.TradeRecharge)
	server.POST("/trade/order_list", handler.TradeOrderList)
	server.POST("/trade/order", handler.TradeOrder)

	server.POST("/tag/list_scene", handler.GetTagListByScene)
	server.POST("/tag/list_item", handler.GetTagListByItem)
	server.POST("/tag/list_user", handler.GetTagListByUser)
	server.POST("/tag/user_interest", handler.UpdateUserTags)

	server.POST("/event/tracking_report", handler.TrackingReport)

	server.POST("/img/upload", handler.PicHostUpload)
	server.GET("/img/download/*url", handler.PicHostDownload)
}
