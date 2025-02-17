package infos

import (
	"discord_drive/list"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInfos(c *gin.Context) {
	filename := c.Query("filename")
	fileChannels, err := list.ListChannelFile(c)
	if err != nil {
		return
	}

	if !list.ContainChannel(fileChannels, filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le fichier (" + filename + ") n'existe pas"})
		return
	}

	channel := list.GetChannel(fileChannels, filename)

	c.JSON(http.StatusOK, channel)
}
