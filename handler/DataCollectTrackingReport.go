package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/data_collector"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/jlu-cow-studio/common/model/http_struct/event"
)

func TrackingReport(c *gin.Context) {
	trackingReportReq := new(event.TrackingReportReq)
	trackingReportRes := &event.TrackingReportRes{
		Base: http_struct.ResBase{
			Code:    "400",
			Message: "failed",
		},
	}
	defer c.JSON(200, trackingReportRes)
	defer log.Println("response: ", trackingReportRes)

	bodyb, ok := c.Get("body")
	if !ok {
		trackingReportRes.Base.Message = "error when get body"
		return
	}

	log.Println("request body: ", string(bodyb.([]byte)))
	if err := json.Unmarshal(bodyb.([]byte), trackingReportReq); err != nil {
		trackingReportRes.Base.Message = err.Error()
		return
	}

	log.Println("request: ", trackingReportReq)

	cli, err := rpc.GetDataCollectorCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		trackingReportRes.Base.Message = err.Error()
		return
	}

	rpcTrackingReportReq := &data_collector.TrackingReportReq{
		Base: &base.BaseReq{
			Token: trackingReportReq.Base.Token,
			Logid: trackingReportReq.Base.LogId,
		},
		ItemId:   trackingReportReq.ItemID,
		Behavior: trackingReportReq.Behavior,
	}

	log.Printf("rpc request: %+v\n", rpcTrackingReportReq)

	rpcTrackingReportRes, err := cli.TrackingReport(context.Background(), rpcTrackingReportReq)
	log.Printf("rpc response: %+v\n", rpcTrackingReportRes)

	if err != nil {
		log.Printf("rpc call error: %s\n", err.Error())
		trackingReportRes.Base.Message = err.Error()
		return
	}

	trackingReportRes.Base.Code = rpcTrackingReportRes.Base.Code
	trackingReportRes.Base.Message = rpcTrackingReportRes.Base.Message
}
