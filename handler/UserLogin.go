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
	"github.com/sanity-io/litter"
)

func UserLogin(c *gin.Context) {
	loginReq := new(user.UserLoginReq)
	loginRes := &user.UserLoginRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, loginRes)
	defer log.Println("response: ", litter.Sdump(loginRes))

	bodyb, ok := c.Get("body")
	if !ok {
		loginRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), loginReq); err != nil {
		loginRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", litter.Sdump(loginReq))

	cli, err := rpc.GetUserCoreCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		loginRes.Base.Message = err.Error()
		return
	}

	rpcUserLoginReq := &user_core.UserLoginReq{
		Base: &base.BaseReq{
			Token: loginReq.Base.Token,
			Logid: loginReq.Base.LogId,
		},
		Username: loginReq.Username,
		Password: loginReq.Password,
	}

	log.Printf("rpc request: %+v\n", rpcUserLoginReq)

	rpcUserLoginRes, err := cli.UserLogin(context.Background(), rpcUserLoginReq)
	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		loginRes.Base.Message = err.Error()
		return
	}

	loginRes.Base.Code = rpcUserLoginRes.Base.Code
	loginRes.Base.Message = rpcUserLoginRes.Base.Message
	loginRes.Token = rpcUserLoginRes.Token
	log.Printf("rpc response: %+v\n", rpcUserLoginRes)
}
