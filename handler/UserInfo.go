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

func UserInfo(c *gin.Context) {
	userInfoReq := new(http_struct.ReqBase)
	userInfoRes := &user.UserInfoRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, userInfoRes)
	defer log.Println("response: ", userInfoRes)

	bodyb, ok := c.Get("body")
	if !ok {
		userInfoRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), userInfoReq); err != nil {
		userInfoRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", userInfoReq)

	conn, err := rpc.GetConn(UserCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		userInfoRes.Base.Message = err.Error()
		return
	}

	cli := user_core.NewUserCoreServiceClient(conn)

	rpcUserInfoReq := &user_core.UserInfoReq{
		Base: &base.BaseReq{
			Token: userInfoReq.Token,
			Logid: userInfoReq.LogId,
		},
	}

	log.Printf("rpc request: %+v\n", rpcUserInfoReq)

	rpcUserInfoRes, err := cli.UserInfo(context.Background(), rpcUserInfoReq)
	log.Printf("rpc response: %+v\n", rpcUserInfoRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		userInfoRes.Base.Message = err.Error()
		return
	}

	userInfoRes.Base.Code = rpcUserInfoRes.Base.Code
	userInfoRes.Base.Message = rpcUserInfoRes.Base.Message
	userInfoRes.Username = rpcUserInfoRes.UserInfo.Username
	userInfoRes.Province = rpcUserInfoRes.UserInfo.Province
	userInfoRes.City = rpcUserInfoRes.UserInfo.City
	userInfoRes.District = rpcUserInfoRes.UserInfo.District
	userInfoRes.Role = rpcUserInfoRes.UserInfo.Role
}
