package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/tag_core"
	"github.com/jlu-cow-studio/common/model/http_struct"
)

func GetTagListByUser(c *gin.Context) {
	getTagListByUserReq := new(tag.GetTagListByUserRequest)
	getTagListByUserRes := &tag.GetTagListByUserResponse{
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

	conn, err := rpc.GetConn(TagCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		getTagListByUserRes.Base.Message = err.Error()
		return
	}

	cli := tag_core.NewTagCoreServiceClient(conn)

	rpcGetTagListByUserReq := &tag_core.GetTagListByUserRequest{
		Base: &base.BaseReq{
			Token: getTagListByUserReq.Base.Token,
			Logid: getTagListByUserReq.Base.LogId,
		},
		UserId: getTagListByUserReq.UserId,
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
	for _, tcl := range rpcGetTagListByUserRes.List {
		tcwl := new(tag.TagCategoryWithList)
		tcwl.TagCategory.ID = int(tcl.Category.Id)
		tcwl.TagCategory.Name = tcl.Category.Name
		tcwl.TagCategory.ParentID = int(tcl.Category.ParentId)
		tcwl.TagCategory.Level = tcl.Category.Level
		for _, tl := range tcl.TagList {
			t := new(tag.Tag)
			t.ID = tl.Id
			t.Name = tl.Name
			t.Weight = tl.Weight
			t.MarkObject = tl.MarkObject
			t.CategoryID = tl.CategoryId
			tcwl.TagList = append(tcwl.TagList, t)
		}
		getTagListByUserRes.List = append(getTagListByUserRes.List, tcwl)
	}
}
