package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/user_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/user"
)

func UserFollow(c *gin.Context) {
	followReq := new(user.UserFollowReq)
	followRes := &user.UserFollowRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, followRes)
	defer log.Println("response: ", followRes)

	bodyb, ok := c.Get("body")
	if !ok {
		followRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), followReq); err != nil {
		followRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", followReq)

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(http_struct.UserTokenInfo)

	if tokenInfo.LoginState != http_struct.LoggedIn {
		followRes.Base.Message = "please login first"
		followRes.Base.Code = "401"
		return
	}

	conn, err := rpc.GetConn(UserCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		followRes.Base.Message = err.Error()
		return
	}

	cli := user_core.NewUserCoreServiceClient(conn)

	rpcFollowReq := &user_core.FollowReq{
		Base: &base.BaseReq{
			Token: followReq.Base.Token,
			Logid: followReq.Base.LogId,
		},
		FollowerId:  tokenInfo.Uid,
		FollowingId: followReq.FollowingId,
		Action:      followReq.Action,
	}

	log.Printf("rpc request: %+v\n", rpcFollowReq)

	rpcFollowRes, err := cli.Follow(context.Background(), rpcFollowReq)
	log.Printf("rpc response: %+v\n", rpcFollowRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		followRes.Base.Message = err.Error()
		return
	}

	followRes.Base.Code = rpcFollowRes.Base.Code
	followRes.Base.Message = rpcFollowRes.Base.Message
}
