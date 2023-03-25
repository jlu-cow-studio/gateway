package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/item"
)

func ItemAdd(c *gin.Context) {
	addItemReq := new(item.AddItemReq)
	addItemRes := &item.AddItemRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, addItemRes)
	defer log.Println("response: ", addItemRes)

	bodyb, ok := c.Get("body")
	if !ok {
		addItemRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), addItemReq); err != nil {
		addItemRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", addItemReq)

	conn, err := rpc.GetConn(ProductCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		addItemRes.Base.Message = err.Error()
		return
	}

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(*http_struct.UserTokenInfo)

	uid, err := strconv.Atoi(tokenInfo.Uid)
	if err != nil {
		log.Printf("get uid error: %s\n", err.Error())
		addItemRes.Base.Message = err.Error()
		return
	}

	cli := product_core.NewProductCoreServiceClient(conn)

	rpcAddItemReq := &product_core.AddItemReq{
		Base: &base.BaseReq{
			Token: addItemReq.Base.Token,
			Logid: addItemReq.Base.LogId,
		},
		ItemInfo: &product_core.ItemInfo{
			Name:               addItemReq.Name,
			Description:        addItemReq.Description,
			Category:           addItemReq.Category,
			Price:              addItemReq.Price,
			Stock:              addItemReq.Stock,
			Province:           addItemReq.Province,
			City:               addItemReq.City,
			District:           addItemReq.District,
			UserId:             int32(uid),
			UserType:           tokenInfo.Role,
			SpecificAttributes: addItemReq.SpecificAttr,
		},
	}

	log.Printf("rpc request: %+v\n", rpcAddItemReq)

	rpcAddItemRes, err := cli.AddItem(context.Background(), rpcAddItemReq)
	log.Printf("rpc response: %+v\n", rpcAddItemRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		addItemRes.Base.Message = err.Error()
		return
	}

	addItemRes.Base.Code = rpcAddItemRes.Base.Code
	addItemRes.Base.Message = rpcAddItemRes.Base.Message
	addItemRes.ItemID = rpcAddItemRes.GetItemId()
}
