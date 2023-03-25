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

func ItemDelete(c *gin.Context) {
	deleteItemReq := new(item.DeleteItemReq)
	deleteItemRes := &item.DeleteItemRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, deleteItemRes)
	defer log.Println("response: ", deleteItemRes)

	bodyb, ok := c.Get("body")
	if !ok {
		deleteItemRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), deleteItemReq); err != nil {
		deleteItemRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", deleteItemReq)

	conn, err := rpc.GetConn(ProductCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		deleteItemRes.Base.Message = err.Error()
		return
	}

	cli := product_core.NewProductCoreServiceClient(conn)

	rpcDeleteItemReq := &product_core.DeleteItemReq{
		Base: &base.BaseReq{
			Token: deleteItemReq.Base.Token,
			Logid: deleteItemReq.Base.LogId,
		},
		ItemId: deleteItemReq.ItemID,
	}

	log.Printf("rpc request: %+v\n", rpcDeleteItemReq)

	rpcDeleteItemRes, err := cli.DeleteItem(context.Background(), rpcDeleteItemReq)
	log.Printf("rpc response: %+v\n", rpcDeleteItemRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		deleteItemRes.Base.Message = err.Error()
		return
	}

	deleteItemRes.Base.Code = rpcDeleteItemRes.Base.Code
	deleteItemRes.Base.Message = rpcDeleteItemRes.Base.Message
}
