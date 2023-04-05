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

func GetTagListByUser(c *gin.Context) {
	getTagListByUserReq := new(tag.GetTagListByUserReq)
	getTagListByUserRes := &tag.GetTagListByUserRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, getTagListByUserRes)
	defer log.Println("response: ", getTagListByUserRes)

	bodyb, ok := c.Get("body")
	if !ok {
		getTagListByUserRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), getTagListByUserReq); err != nil {
		getTagListByUserRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", getTagListByUserReq)

	cli, err := rpc.GetTagCoreCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		getTagListByUserRes.Base.Message = err.Error()
		return
	}

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(*http_struct.UserTokenInfo)

	rpcGetTagListByUserReq := &tag_core.GetTagListByUserRequest{
		Base: &base.BaseReq{
			Token: getTagListByUserReq.Base.Token,
			Logid: getTagListByUserReq.Base.LogId,
		},
		UserId: tokenInfo.Uid,
	}

	log.Printf("rpc request: %+v\n", rpcGetTagListByUserReq)

	rpcGetTagListByUserRes, err := cli.GetTagListByUser(context.Background(), rpcGetTagListByUserReq)
	log.Printf("rpc response: %+v\n", rpcGetTagListByUserRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		getTagListByUserRes.Base.Message = err.Error()
		return
	}

	getTagListByUserRes.Base.Code = rpcGetTagListByUserRes.Base.Code
	getTagListByUserRes.Base.Message = rpcGetTagListByUserRes.Base.Message

	// 转换 TagCategoryWithList
	for _, t := range rpcGetTagListByUserRes.TagList {
		tagItem := new(tag.Tag)
		tagItem.ID = t.Id
		tagItem.Name = t.Name
		tagItem.Weight = t.Weight
		tagItem.MarkObject = t.MarkObject
		tagItem.CategoryID = t.CategoryId
		getTagListByUserRes.TagList = append(getTagListByUserRes.TagList, tagItem)
	}
}
