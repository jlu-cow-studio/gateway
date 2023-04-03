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

const TagCoreServiceName = "cowstudio/tag-core"

func GetTagListByScene(c *gin.Context) {
	getTagListBySceneReq := new(tag.GetTagListBySceneRequest)
	getTagListBySceneRes := &tag.GetTagListBySceneResponse{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, getTagListBySceneRes)
	defer log.Println("response: ", getTagListBySceneRes)

	bodyb, ok := c.Get("body")
	if !ok {
		getTagListBySceneRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), getTagListBySceneReq); err != nil {
		getTagListBySceneRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", getTagListBySceneReq)

	conn, err := rpc.GetConn(TagCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		getTagListBySceneRes.Base.Message = err.Error()
		return
	}

	cli := tag_core.NewTagCoreServiceClient(conn)

	rpcGetTagListBySceneReq := &tag_core.GetTagListBySceneRequest{
		Base: &base.BaseReq{
			Token: getTagListBySceneReq.Base.Token,
			Logid: getTagListBySceneReq.Base.LogId,
		},
		Scene: getTagListBySceneReq.Scene,
	}

	log.Printf("rpc request: %+v\n", rpcGetTagListBySceneReq)

	rpcGetTagListBySceneRes, err := cli.GetTagListByScene(context.Background(), rpcGetTagListBySceneReq)
	log.Printf("rpc response: %+v\n", rpcGetTagListBySceneRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		getTagListBySceneRes.Base.Message = err.Error()
		return
	}

	getTagListBySceneRes.Base.Code = rpcGetTagListBySceneRes.Base.Code
	getTagListBySceneRes.Base.Message = rpcGetTagListBySceneRes.Base.Message

	// 转换 TagCategoryWithList
	for _, tcl := range rpcGetTagListBySceneRes.List {
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
		getTagListBySceneRes.List = append(getTagListBySceneRes.List, tcwl)
	}
}
