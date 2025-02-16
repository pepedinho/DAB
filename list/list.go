package list

import (
	"discord_drive/upload"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func ListFile(c *gin.Context) {
	dg, err := discordgo.New("Bot " + upload.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'instancier le bot discord : " + err.Error()})
		return
	}

	err = dg.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de se connecter a discord : " + err.Error()})
		return
	}

	defer dg.Close() // fermer avant de return

	list, err := dg.GuildChannels(upload.GuildID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de récuperé les fichiers : " + err.Error()})
		return
	}

	var fileChannels []map[string]interface{}
	for _, channel := range list {
		if channel.Type == discordgo.ChannelTypeGuildText {
			parts := strings.SplitN(channel.Name, "__", 2)
			if len(parts) == 2 {
				fileChannels = append(fileChannels, map[string]interface{}{
					"id":        channel.ID,
					"channel":   channel.Name,
					"file_name": parts[1],
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"files": fileChannels})
}
