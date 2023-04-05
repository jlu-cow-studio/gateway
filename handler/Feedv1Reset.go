package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/feed"
)

func Feedv1Reset(c *gin.Context) {
	resetFeedReq := new(feed.ResetFeedReq)
	resetFeedRes := &feed.ResetFeedRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, resetFeedRes)
	defer log.Println("response: ", resetFeedRes)

	bodyb, ok := c.Get("body")
	if !ok {
		resetFeedRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), resetFeedReq); err != nil {
		resetFeedRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", resetFeedReq)

	cli, err := rpc.GetFeedCli()
	if err != nil {
		log.Printf("get rpc cli error: %s\n", err.Error())
		resetFeedRes.Base.Message = err.Error()
		return
	}

	rpcResetFeedReq := &feed_service.ResetFeedRequest{
		Base: &base.BaseReq{
			Token: resetFeedReq.Base.Token,
			Logid: resetFeedReq.Base.LogId,
		},
		Scene: resetFeedReq.Scene,
	}

	log.Printf("rpc request: %+v\n", rpcResetFeedReq)

	rpcResetFeedRes, err := cli.ResetFeed(context.Background(), rpcResetFeedReq)
	log.Printf("rpc response: %+v\n", rpcResetFeedRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		resetFeedRes.Base.Message = err.Error()
		return
	}

	resetFeedRes.Base.Code = rpcResetFeedRes.Base.Code
	resetFeedRes.Base.Message = rpcResetFeedRes.Base.Message
}
