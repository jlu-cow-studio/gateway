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

func UserRegister(c *gin.Context) {
	regReq := new(user.UserRegisterReq)
	regRes := &user.UserRegisterRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, regRes)
	defer log.Println("response: ", litter.Sdump(regRes))

	bodyb, ok := c.Get("body")
	if !ok {
		regRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), regReq); err != nil {
		regRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", litter.Sdump(regReq))

	cli, err := rpc.GetUserCoreCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		regRes.Base.Message = err.Error()
		return
	}

	rpcUserRegisterReq := &user_core.UserRegisterReq{
		Base: &base.BaseReq{
			Token: regReq.Base.Token,
			Logid: regReq.Base.LogId,
		},
		UserInfo: &user_core.UserInfo{
			Username: regReq.Username,
			Password: regReq.Password,
			Province: regReq.Province,
			City:     regReq.City,
			District: regReq.District,
		},
	}

	log.Printf("rpc request: %+v\n", rpcUserRegisterReq)

	rpcUserRegisterRes, err := cli.UserRegister(context.Background(), rpcUserRegisterReq)
	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		regRes.Base.Message = err.Error()
		return
	}

	regRes.Base.Code = rpcUserRegisterRes.Base.Code
	regRes.Base.Message = rpcUserRegisterRes.Base.Message
	log.Printf("rpc response: %+v\n", rpcUserRegisterRes)
}
