package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/trade_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/trade"
)

const TradeCoreServiceName = "cowstudio/trade-core"

func TradeRecharge(c *gin.Context) {
	rechargeReq := new(trade.RechargeReq)
	rechargeRes := &trade.RechargeRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(http.StatusOK, rechargeRes)
	defer log.Println("response: ", rechargeRes)

	bodyb, ok := c.Get("body")
	if !ok {
		rechargeRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), rechargeReq); err != nil {
		rechargeRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", rechargeReq)

	conn, err := rpc.GetConn(TradeCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		rechargeRes.Base.Message = err.Error()
		return
	}

	cli := trade_core.NewTradeCoreServiceClient(conn)

	rpcRechargeReq := &trade_core.RechargeRequest{
		Base: &base.BaseReq{
			Token: rechargeReq.Base.Token,
			Logid: rechargeReq.Base.LogId,
		},
		Money: rechargeReq.Money,
	}

	log.Printf("rpc request: %+v\n", rpcRechargeReq)

	rpcRechargeRes, err := cli.Recharge(context.Background(), rpcRechargeReq)
	log.Printf("rpc response: %+v\n", rpcRechargeRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		rechargeRes.Base.Message = err.Error()
		return
	}

	rechargeRes.Base.Code = rpcRechargeRes.Base.Code
	rechargeRes.Base.Message = rpcRechargeRes.Base.Message
}
