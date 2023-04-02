package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/trade_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/item"
	"github.com/jlu-cow-studio/common/model/http_struct/trade"
)

func TradeOrderList(c *gin.Context) {
	orderListReq := new(trade.OrderListReq)
	orderListRes := &trade.OrderListRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, orderListRes)
	defer log.Println("response: ", orderListRes)

	bodyb, ok := c.Get("body")
	if !ok {
		orderListRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), orderListReq); err != nil {
		orderListRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", orderListReq)

	conn, err := rpc.GetConn(TradeCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		orderListRes.Base.Message = err.Error()
		return
	}

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(*http_struct.UserTokenInfo)

	cli := trade_core.NewTreadeCoreServiceClient(conn)

	uid, err := strconv.Atoi(tokenInfo.Uid)
	if err != nil {
		log.Printf("get uid error: %s\n", err.Error())
		orderListRes.Base.Message = err.Error()
		return
	}
	rpcOrderListReq := &trade_core.OrderListRequest{
		Base: &base.BaseReq{
			Token: orderListReq.Base.Token,
			Logid: orderListReq.Base.LogId,
		},
		UserId:   int64(uid),
		Page:     orderListReq.Page,
		PageSize: orderListReq.PerPage,
	}

	log.Printf("rpc request: %+v\n", rpcOrderListReq)

	rpcOrderListRes, err := cli.OrderList(context.Background(), rpcOrderListReq)
	log.Printf("rpc response: %+v\n", rpcOrderListRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		orderListRes.Base.Message = err.Error()
		return
	}

	orderListRes.Base.Code = rpcOrderListRes.Base.Code
	orderListRes.Base.Message = rpcOrderListRes.Base.Message
	orderListRes.OrderList = make([]*trade.OrderInfo, len(rpcOrderListRes.OrderList))
	for i, o := range rpcOrderListRes.OrderList {
		orderListRes.OrderList[i] = &trade.OrderInfo{
			ID:        o.Id,
			UserID:    o.UserId,
			ItemID:    o.ItemId,
			Quantity:  o.Quantity,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			ItemInfo: &item.ItemInfo{
				ID:                 o.ItemInfo.ItemId,
				Name:               o.ItemInfo.Name,
				Description:        o.ItemInfo.Description,
				Category:           o.ItemInfo.Category,
				Price:              o.ItemInfo.Price,
				Stock:              o.ItemInfo.Stock,
				ImageURL:           o.ItemInfo.ImageUrl,
				Province:           o.ItemInfo.Province,
				City:               o.ItemInfo.City,
				District:           o.ItemInfo.District,
				UserID:             o.ItemInfo.UserId,
				UserType:           o.ItemInfo.UserType,
				SpecificAttributes: o.ItemInfo.SpecificAttributes,
			},
		}
	}

	log.Println("response: ", orderListRes)
}
