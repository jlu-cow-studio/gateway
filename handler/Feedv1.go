package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/feed"
)

func Feedv1(c *gin.Context) {
	getFeedReq := new(feed.GetFeedReq)
	getFeedRes := &feed.GetFeedRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}

	defer c.JSON(200, getFeedRes)

	bodyb, ok := c.Get("body")
	if !ok {
		getFeedRes.Base.Message = "error when getting body"
		return
	}

	log.Println("Request body:", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), getFeedReq); err != nil {
		getFeedRes.Base.Message = err.Error()
		return
	}

	log.Println("Request:", getFeedReq)

	cli, err := rpc.GetFeedCli()
	if err != nil {
		log.Printf("Error when getting RPC cli: %s\n", err.Error())
		getFeedRes.Base.Message = err.Error()
		return
	}

	rpcGetFeedReq := &feed_service.GetFeedRequest{
		Base: &base.BaseReq{
			Token: getFeedReq.Base.Token,
			Logid: getFeedReq.Base.LogId,
		},
		Scene:    getFeedReq.Scene,
		Page:     getFeedReq.Page,
		PageSize: getFeedReq.PageSize,
	}

	log.Printf("RPC request: %+v\n", rpcGetFeedReq)

	rpcGetFeedRes, err := cli.GetFeed(context.Background(), rpcGetFeedReq)
	log.Printf("RPC response: %+v\n", rpcGetFeedRes)

	if err != nil {
		log.Printf("Error when making RPC call: %s\n", err.Error())
		getFeedRes.Base.Message = err.Error()
		return
	}

	getFeedRes.Base.Code = rpcGetFeedRes.Base.Code
	getFeedRes.Base.Message = rpcGetFeedRes.Base.Message
	getFeedRes.Items = make([]*feed.ItemForFeed, len(rpcGetFeedRes.Items))
	for i, info := range rpcGetFeedRes.Items {
		getFeedRes.Items[i] = &feed.ItemForFeed{
			ID:                 int(info.Id),
			Name:               info.Name,
			Description:        info.Description,
			Category:           info.Category,
			Price:              info.Price,
			Stock:              int(info.Stock),
			Province:           info.Province,
			City:               info.City,
			District:           info.District,
			ImageURL:           info.ImageUrl,
			UserID:             int(info.UserId),
			UserType:           info.UserType,
			SpecificAttributes: info.SpecificAttributes,
			UID:                int(info.Uid),
			Username:           info.Username,
			UProvince:          info.Uprovince,
			UCity:              info.Ucity,
			UDistrict:          info.Udistrict,
			URole:              info.Urole,
		}
	}
	getFeedRes.Page = rpcGetFeedRes.Pagination.CurrentPage
	getFeedRes.PageSize = rpcGetFeedRes.Pagination.ItemsPerPage
	getFeedRes.Total = rpcGetFeedRes.Pagination.TotalItems
}
