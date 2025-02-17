package main

import (
	"discord_drive/srcs/list"
	"discord_drive/srcs/upload"
	"discord_drive/srcs/get"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("/upload", upload.UploadFile)
	r.GET("/list", list.ListFile)
	r.GET("/get", get.GetFile)

	r.Run(":8000")
}
