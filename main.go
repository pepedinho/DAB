package main

import (
	"discord_drive/upload"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("/upload", upload.UploadFile)

	r.Run(":8000")
}
