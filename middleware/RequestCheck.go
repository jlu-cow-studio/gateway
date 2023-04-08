package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jlu-cow-studio/common/model/http_struct"
)

func RequestCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/img/") {
			return
		}

		baseReq := &http_struct.OnlyBaseReq{}

		buf := &bytes.Buffer{}
		buf.ReadFrom(ctx.Request.Body)

		fmt.Println("checking request: ", buf.String())
		ctx.Set("body", buf.Bytes())

		if err := json.Unmarshal(buf.Bytes(), baseReq); err != nil {
			ctx.AbortWithStatusJSON(200, http_struct.InvalidRequest)
			return
		}

		ctx.Set("base", baseReq.Base)
	}
}
