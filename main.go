package main

import (
	"discord_drive/common"
	"discord_drive/get"
	"discord_drive/infos"
	"discord_drive/list"
	"discord_drive/upload"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	// originsString := "https://dab-frontend-b5b1.vercel.app, http://localhost:3000, http://10.255.255.254:3000"
	// var allowedOrigins []string
	// if originsString != "" {
	// 	allowedOrigins = strings.Split(originsString, ",")
	// }

	return func(c *gin.Context) {
		// isOriginAllowed := func(origin string, allowedOrigins []string) bool {
		// 	for _, allowedOrigin := range allowedOrigins {
		// 		if origin == allowedOrigin {
		// 			return true
		// 		}
		// 	}
		// 	return false
		// }

		origin := c.Request.Header.Get("Origin")

		if origin == "" {
			origin = "*"
		}

		gin.DefaultWriter.Write([]byte("📌 Writing Header\n"))
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// if isOriginAllowed(origin, allowedOrigins) {
		// }

		// Handle preflight OPTIONS requests by aborting with status 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Call the next handler
		c.Next()
	}
}

func main() {

	r := gin.Default()

	r.Use(corsMiddleware())

	r.POST("/upload/:guildID", upload.UploadFile)
	r.GET("/list/:guildID", list.ListFile)
	r.GET("/get/:guildID", get.GetFile)
	r.GET("/infos/:guildID", infos.GetInfos)

	common.InitDiscordSession()
	defer common.DiscordSession.Close()

	r.Run("0.0.0.0:8000")
}
