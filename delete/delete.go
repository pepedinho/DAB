package delete

import (
	"discord_drive/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteFile(c *gin.Context) {
	dg := common.DiscordSession

	channelID := c.Param("channelID")

	_, err := dg.ChannelDelete(channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer le fichier : " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Fichier supprimé avec succés"})
}
