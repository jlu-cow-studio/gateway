package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const BasePath = "/opt/img/"
const BaseUrl = "/img/download"

// 获取随机的文件名
func getRandomName() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

// 处理上传请求
func PicHostUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No file is uploaded"})
		return
	}

	// 遍历文件列表，依次保存文件
	var urls []string
	for _, file := range files {
		fileName := getRandomName() + filepath.Ext(file.Filename)
		filePath := BasePath + fileName

		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
			return
		}

		url := BaseUrl + fileName
		urls = append(urls, url)
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
