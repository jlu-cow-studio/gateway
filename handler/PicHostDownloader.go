package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 处理下载请求
func PicHostDownload(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing url parameter"})
		return
	}

	filePath := "/opt/img/" + filepath.Base(url)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"message": "File not found"})
		return
	}

	// 读取文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read file"})
		return
	}

	// 发送文件
	c.Data(http.StatusOK, "image/png", data)
}
