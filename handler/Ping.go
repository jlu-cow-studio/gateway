package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	fmt.Fprintln(c.Writer, "pong")
}
