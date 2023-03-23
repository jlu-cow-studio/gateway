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

func UserAuth(c *gin.Context) {
	authReq := new(user.UserAuthReq)
	authRes := &http_struct.ResBase{
		Code:    "400",
		Message: "failed",
	}
	defer c.JSON(200, authRes)
	defer log.Println("response: ", litter.Sdump(authRes))

	bodyb, ok := c.Get("body")
	if !ok {
		authRes.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), authReq); err != nil {
		authRes.Message = err.Error()
		return
	}

	log.Println("request: ", litter.Sdump(authReq))

	conn, err := rpc.GetConn(UserCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		authRes.Message = err.Error()
		return
	}

	cli := user_core.NewUserCoreServiceClient(conn)

	rpcUserAuthReq := &user_core.UserAuthReq{
		Base: &base.BaseReq{
			Token: authReq.Base.Token,
			Logid: authReq.Base.LogId,
		},
		Role: authReq.Role,
	}

	log.Printf("rpc request: %+v\n", rpcUserAuthReq)

	rpcUserAuthRes, err := cli.UserAuth(context.Background(), rpcUserAuthReq)
	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		authRes.Message = err.Error()
		return
	}

	authRes.Code = rpcUserAuthRes.Base.Code
	authRes.Message = rpcUserAuthRes.Base.Message
	log.Printf("rpc response: %+v\n", rpcUserAuthRes)
}
