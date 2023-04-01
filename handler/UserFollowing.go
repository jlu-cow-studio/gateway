package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/user_core"
	"github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/user"
)

func UserFollowing(c *gin.Context) {
	followingReq := new(user.UserFollowingReq)
	followingRes := &user.UserFollowingRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, followingRes)
	defer log.Println("response: ", followingRes)

	bodyb, ok := c.Get("body")
	if !ok {
		followingRes.Base.Message = "error when getting request body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), followingReq); err != nil {
		followingRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", followingReq)

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(http_struct.UserTokenInfo)

	if tokenInfo.LoginState != http_struct.LoggedIn {
		followingRes.Base.Message = "please login first"
		followingRes.Base.Code = "401"
		return
	}

	conn, err := rpc.GetConn(UserCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		followingRes.Base.Message = err.Error()
		return
	}

	cli := user_core.NewUserCoreServiceClient(conn)

	rpcFollowingReq := &user_core.FollowingReq{
		Base: &base.BaseReq{
			Token: tokenInfo.Token,
			Logid: followingReq.Base.LogId,
		},
		UserId:   tokenInfo.Uid,
		Page:     int32(followingReq.Page),
		PageSize: int32(followingReq.PageSize),
	}

	log.Printf("rpc request: %+v\n", rpcFollowingReq)

	rpcFollowingRes, err := cli.Following(context.Background(), rpcFollowingReq)
	log.Printf("rpc response: %+v\n", rpcFollowingRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		followingRes.Base.Message = err.Error()
		return
	}

	followingRes.Base.Code = rpcFollowingRes.Base.Code
	followingRes.Base.Message = rpcFollowingRes.Base.Message
	followingRes.TotalCount = int(rpcFollowingRes.TotalCount)
	followingRes.TotalPage = int(rpcFollowingRes.TotalPage)
	followingRes.Users = make([]*redis.UserInfo, len(rpcFollowingRes.Users))

	for i, userInfo := range rpcFollowingRes.Users {
		followingRes.Users[i] = &redis.UserInfo{
			Uid:      userInfo.Uid,
			Username: userInfo.Username,
			Role:     userInfo.Role,
		}
	}
}
