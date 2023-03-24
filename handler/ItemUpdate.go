package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/item"
)

const ProductCoreServiceName = "cowstudio/product-core"

func ItemUpdate(c *gin.Context) {
	updateItemReq := new(item.UpdateItemReq)
	updateItemRes := &item.UpdateItemRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, updateItemRes)
	defer log.Println("response: ", updateItemRes)

	bodyb, ok := c.Get("body")
	if !ok {
		updateItemRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), updateItemReq); err != nil {
		updateItemRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", updateItemReq)

	conn, err := rpc.GetConn(ProductCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		updateItemRes.Base.Message = err.Error()
		return
	}

	cli := product_core.NewProductCoreServiceClient(conn)

	rpcUpdateItemReq := &product_core.UpdateItemReq{
		Base: &base.BaseReq{
			Token: updateItemReq.Base.Token,
			Logid: updateItemReq.Base.LogId,
		},
		Item: &product_core.ItemInfo{
			ItemId:             updateItemReq.ItemID,
			Name:               updateItemReq.Name,
			Description:        updateItemReq.Description,
			Price:              updateItemReq.Price,
			Stock:              updateItemReq.Stock,
			Province:           updateItemReq.Province,
			City:               updateItemReq.City,
			District:           updateItemReq.District,
			SpecificAttributes: updateItemReq.SpecificAttr,
		},
	}

	log.Printf("rpc request: %+v\n", rpcUpdateItemReq)

	rpcUpdateItemRes, err := cli.UpdateItem(context.Background(), rpcUpdateItemReq)
	log.Printf("rpc response: %+v\n", rpcUpdateItemRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		updateItemRes.Base.Message = err.Error()
		return
	}

	updateItemRes.Base.Code = rpcUpdateItemRes.Base.Code
	updateItemRes.Base.Message = rpcUpdateItemRes.Base.Message
}
