package handler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/user_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/user"
)

const UserCoreServiceName = "cowstudio/user-core"

func UserLogin(c *gin.Context) {

	loginReq := new(user.UserLoginReq)
	loginRes := &user.UserLoginRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, loginRes)

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(c.Request.Body); err != nil {
		loginRes.Base.Message = err.Error()
		return
	}

	if err := json.Unmarshal(buf.Bytes(), loginReq); err != nil {
		loginRes.Base.Message = err.Error()
		return
	}

	conn, err := rpc.GetConn(UserCoreServiceName)
	if err != nil {
		loginRes.Base.Message = err.Error()
		return
	}

	cli := user_core.NewUserCoreServiceClient(conn)

	rpcUserLoginReq := &user_core.UserLoginReq{
		Base: &base.BaseReq{
			Token: loginReq.Base.Token,
			Logid: loginReq.Base.LogId,
		},
		Username: loginReq.Username,
		Password: loginReq.Password,
	}

	rpcUserLoginRes, _ := cli.UserLogin(context.Background(), rpcUserLoginReq)

	loginRes.Base.Code = rpcUserLoginRes.Base.Code
	loginRes.Base.Message = rpcUserLoginRes.Base.Message
	loginRes.Token = rpcUserLoginRes.Token
}
