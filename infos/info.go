package infos

import (
	"discord_drive/common"
	"discord_drive/list"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInfos(c *gin.Context) {
	filename := c.Query("filename")
	dg := common.DiscordSession
	fileChannels, err := list.ListChannelFileWithDg(c, dg)
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
