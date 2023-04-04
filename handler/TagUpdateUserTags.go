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

func UpdateUserTags(c *gin.Context) {
	updateUserTagsReq := new(tag.UpdateUserTagsReq)
	updateUserTagsRes := &tag.UpdateUserTagsRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, updateUserTagsRes)
	defer log.Println("response: ", updateUserTagsRes)

	bodyb, ok := c.Get("body")
	if !ok {
		updateUserTagsRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), updateUserTagsReq); err != nil {
		updateUserTagsRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", updateUserTagsReq)

	conn, err := rpc.GetConn(TagCoreServiceName)
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		updateUserTagsRes.Base.Message = err.Error()
		return
	}

	cli := tag_core.NewTagCoreServiceClient(conn)

	tokenInfor, _ := c.Get("tokenInfo")
	tokenInfo := tokenInfor.(*http_struct.UserTokenInfo)

	rpcUpdateUserTagsReq := &tag_core.UpdateUserTagsRequest{
		Base: &base.BaseReq{
			Token: updateUserTagsReq.Base.Token,
			Logid: updateUserTagsReq.Base.LogId,
		},
		TagList: updateUserTagsReq.TagList,
		UserId:  tokenInfo.Uid,
	}

	log.Printf("rpc request: %+v\n", rpcUpdateUserTagsReq)

	rpcUpdateUserTagsRes, err := cli.UpdateUserTags(context.Background(), rpcUpdateUserTagsReq)
	log.Printf("rpc response: %+v\n", rpcUpdateUserTagsRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		updateUserTagsRes.Base.Message = err.Error()
		return
	}

	updateUserTagsRes.Base.Code = rpcUpdateUserTagsRes.Base.Code
	updateUserTagsRes.Base.Message = rpcUpdateUserTagsRes.Base.Message
}
