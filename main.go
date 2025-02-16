package main

import (
	"discord_drive/list"
	"discord_drive/upload"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("/upload", upload.UploadFile)

	r.GET("/list", list.ListFile)

	r.Run(":8000")
}
