package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/dal/redis"
	model_redis "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/common/model/http_struct"
	"github.com/sanity-io/litter"
)

func GetUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmp, _ := ctx.Get("base")

		var base http_struct.ReqBase = tmp.(http_struct.ReqBase)

		tokenInfo := &http_struct.UserTokenInfo{
			Token:      base.Token,
			LoginState: http_struct.NotLogged,
		}

		fmt.Println("checking user token", tokenInfo.Token)

		if cmd := redis.DB.Get(redis.GetUserTokenKey(tokenInfo.Token)); cmd.Err() != nil {
			tokenInfo.LoginState = http_struct.InvalidToken
			fmt.Println("invalid token ", cmd.Err())
		} else {
			infoRaw := cmd.Val()
			info := new(model_redis.UserInfo)
			if err := json.Unmarshal([]byte(infoRaw), info); err != nil {
				tokenInfo.LoginState = http_struct.NotLogged
				fmt.Println("user not logged ", err)
			}

			tokenInfo.Uid = info.Uid
			tokenInfo.Username = info.Username
			tokenInfo.Role = info.Role
			tokenInfo.Province = info.Province
			tokenInfo.City = info.City
			tokenInfo.District = info.District
		}

		fmt.Println("user logged in ", litter.Sdump(tokenInfo))
		ctx.Set("tokenInfo", tokenInfo)
	}
}
