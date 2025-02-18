package main

import (
	"discord_drive/get"
	"discord_drive/infos"
	"discord_drive/list"
	"discord_drive/upload"
	"strings"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	originsString := "http://localhost:3000, http://10.255.255.254:3000, https://dab-frontend-b5b1.vercel.app"
	var allowedOrigins []string
	if originsString != "" {
		allowedOrigins = strings.Split(originsString, ",")
	}

	return func(c *gin.Context) {
		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		}

		origin := c.Request.Header.Get("Origin")

		if isOriginAllowed(origin, allowedOrigins) {
			gin.DefaultWriter.Write([]byte("📌 Writing Header\n"))
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		}

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

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	r.Use(corsMiddleware())

	r.POST("/upload/:guildID", upload.UploadFile)
	r.GET("/list/:guildID", list.ListFile)
	r.GET("/get/:guildID", get.GetFile)
	r.GET("/infos/:guildID", infos.GetInfos)

	r.Run(":8000")
}
