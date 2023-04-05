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

func UserFollowCount(c *gin.Context) {
	followCountReq := new(user.UserFollowCountReq)
	followCountRes := &user.UserFollowCountRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, followCountRes)
	defer log.Println("response: ", followCountRes)

	bodyb, ok := c.Get("body")
	if !ok {
		followCountRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), followCountReq); err != nil {
		followCountRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", followCountReq)

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(http_struct.UserTokenInfo)

	if tokenInfo.LoginState != http_struct.LoggedIn {
		followCountRes.Base.Message = "please login first"
		followCountRes.Base.Code = "401"
		return
	}

	cli, err := rpc.GetUserCoreCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		followCountRes.Base.Message = err.Error()
		return
	}

	rpcFollowCountReq := &user_core.FollowCountReq{
		Base: &base.BaseReq{
			Token: tokenInfo.Token,
			Logid: followCountReq.Base.LogId,
		},
		UserId: tokenInfo.Uid,
	}

	log.Printf("rpc request: %+v\n", rpcFollowCountReq)

	rpcFollowCountRes, err := cli.FollowCount(context.Background(), rpcFollowCountReq)
	log.Printf("rpc response: %+v\n", rpcFollowCountRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		followCountRes.Base.Message = err.Error()
		return
	}

	followCountRes.Base.Code = rpcFollowCountRes.Base.Code
	followCountRes.Base.Message = rpcFollowCountRes.Base.Message
	followCountRes.FollowingCount = int(rpcFollowCountRes.FollowingCount)
	followCountRes.FollowerCount = int(rpcFollowCountRes.FollowerCount)
}
