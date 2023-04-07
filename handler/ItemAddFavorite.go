package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/item"
)

func AddFavorite(c *gin.Context) {
	addFavReq := new(item.AddFavoriteReq)
	addFavRes := &item.AddFavoriteRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(http.StatusOK, addFavRes)
	defer log.Println("response: ", addFavRes)

	bodyb, ok := c.Get("body")
	if !ok {
		addFavRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), addFavReq); err != nil {
		addFavRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", addFavReq)

	cli, err := rpc.GetProductCoreCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		addFavRes.Base.Message = err.Error()
		return
	}

	rpcAddFavReq := &product_core.AddFavoriteReq{
		Base: &base.BaseReq{
			Token: addFavReq.Base.Token,
			Logid: addFavReq.Base.LogId,
		},
		Action: addFavReq.Action,
		ItemId: addFavReq.ItemID,
	}

	log.Printf("rpc request: %+v\n", rpcAddFavReq)

	rpcAddFavRes, err := cli.AddFavorite(context.Background(), rpcAddFavReq)
	log.Printf("rpc response: %+v\n", rpcAddFavRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		addFavRes.Base.Message = err.Error()
		return
	}

	addFavRes.Base.Code = rpcAddFavRes.Base.Code
	addFavRes.Base.Message = rpcAddFavRes.Base.Message
}
