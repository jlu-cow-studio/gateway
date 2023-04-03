package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/tag_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/tag"
)

func GetTagListByItem(c *gin.Context) {
	getTagListByItemReq := new(tag.GetTagListByItemReq)
	getTagListByItemRes := &tag.GetTagListByItemRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, getTagListByItemRes)
	defer log.Println("response: ", getTagListByItemRes)

	bodyb, ok := c.Get("body")
	if !ok {
		getTagListByItemRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), getTagListByItemReq); err != nil {
		getTagListByItemRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", getTagListByItemReq)

	conn, err := rpc.GetConn(TagCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		getTagListByItemRes.Base.Message = err.Error()
		return
	}

	cli := tag_core.NewTagCoreServiceClient(conn)

	rpcGetTagListByItemReq := &tag_core.GetTagListByItemRequest{
		Base: &base.BaseReq{
			Token: getTagListByItemReq.Base.Token,
			Logid: getTagListByItemReq.Base.LogId,
		},
		ItemId: getTagListByItemReq.ItemID,
	}

	log.Printf("rpc request: %+v\n", rpcGetTagListByItemReq)

	rpcGetTagListByItemRes, err := cli.GetTagListByItem(context.Background(), rpcGetTagListByItemReq)
	log.Printf("rpc response: %+v\n", rpcGetTagListByItemRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		getTagListByItemRes.Base.Message = err.Error()
		return
	}

	getTagListByItemRes.Base.Code = rpcGetTagListByItemRes.Base.Code
	getTagListByItemRes.Base.Message = rpcGetTagListByItemRes.Base.Message

	// 转换 Tag
	for _, t := range rpcGetTagListByItemRes.TagList {
		tagItem := new(tag.Tag)
		tagItem.ID = t.Id
		tagItem.Name = t.Name
		tagItem.Weight = t.Weight
		tagItem.MarkObject = t.MarkObject
		tagItem.CategoryID = t.CategoryId
		getTagListByItemRes.TagList = append(getTagListByItemRes.TagList, tagItem)
	}
}
