package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/trade_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/trade"
)

func OrderAdd(c *gin.Context) {
	orderReq := new(trade.OrderReq)
	orderRes := &trade.OrderRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, orderRes)
	defer log.Println("response: ", orderRes)

	bodyb, ok := c.Get("body")
	if !ok {
		orderRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), orderReq); err != nil {
		orderRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", orderReq)

	conn, err := rpc.GetConn(TradeCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		orderRes.Base.Message = err.Error()
		return
	}

	cli := trade_core.NewTreadeCoreServiceClient(conn)

	rpcOrderReq := &trade_core.OrderRequest{
		Base: &base.BaseReq{
			Token: orderReq.Base.Token,
			Logid: orderReq.Base.LogId,
		},
		ItemId: orderReq.ItemID,
		Count:  orderReq.Count,
	}

	log.Printf("rpc request: %+v\n", rpcOrderReq)

	rpcOrderRes, err := cli.Order(context.Background(), rpcOrderReq)
	log.Printf("rpc response: %+v\n", rpcOrderRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		orderRes.Base.Message = err.Error()
		return
	}

	orderRes.Base.Code = rpcOrderRes.Base.Code
	orderRes.Base.Message = rpcOrderRes.Base.Message
}
