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

func UserFollower(c *gin.Context) {
	followerReq := new(user.UserFollowerReq)
	followerRes := &user.UserFollowerRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, followerRes)
	defer log.Println("response: ", followerRes)

	bodyb, ok := c.Get("body")
	if !ok {
		followerRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), followerReq); err != nil {
		followerRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", followerReq)

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(http_struct.UserTokenInfo)

	if tokenInfo.LoginState != http_struct.LoggedIn {
		followerRes.Base.Message = "please login first"
		followerRes.Base.Code = "401"
		return
	}

	conn, err := rpc.GetConn(UserCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		followerRes.Base.Message = err.Error()
		return
	}

	cli := user_core.NewUserCoreServiceClient(conn)

	rpcFollowerReq := &user_core.FollowersReq{
		Base: &base.BaseReq{
			Token: followerReq.Base.Token,
			Logid: followerReq.Base.LogId,
		},
		UserId:   tokenInfo.Uid,
		Page:     int32(followerReq.Page),
		PageSize: int32(followerReq.PageSize),
	}

	log.Printf("rpc request: %+v\n", rpcFollowerReq)

	rpcFollowerRes, err := cli.Followers(context.Background(), rpcFollowerReq)
	log.Printf("rpc response: %+v\n", rpcFollowerRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		followerRes.Base.Message = err.Error()
		return
	}

	followerRes.Base.Code = rpcFollowerRes.Base.Code
	followerRes.Base.Message = rpcFollowerRes.Base.Message
	followerRes.TotalCount = int(rpcFollowerRes.TotalCount)
	followerRes.TotalPage = int(rpcFollowerRes.TotalPage)
	followerRes.Users = make([]*redis.UserInfo, len(rpcFollowerRes.Users))

	for i, userInfo := range rpcFollowerRes.Users {
		followerRes.Users[i] = &redis.UserInfo{
			Uid:      userInfo.Uid,
			Username: userInfo.Username,
			Role:     userInfo.Role,
		}
	}
}
