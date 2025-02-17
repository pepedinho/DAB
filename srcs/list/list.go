package list

import (
	"discord_drive/srcs/common"
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func ListFile(c *gin.Context) {
	fileChannels, err := ListChannelFile(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": fileChannels})
}

func ListChannelFile(c *gin.Context) ([]map[string]interface{}, error) {
	dg, err := discordgo.New("Bot " + common.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'instancier le bot discord : " + err.Error()})
		return nil, fmt.Errorf("impossible d'instancier le bot discord : %w", err)
	}

	err = dg.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de se connecter a discord : " + err.Error()})
		return nil, fmt.Errorf("impossible de se connecter a discord : %w", err)
	}

	defer dg.Close() // fermer avant de return

	list, err := dg.GuildChannels(common.GuildID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de récuperé les fichiers : " + err.Error()})
		return nil, fmt.Errorf("impossible de récupérer les fichiers : %w", err)
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

	return fileChannels, nil
}

// dg has to be already open before call this function
func ListChannelFileWithDg(c *gin.Context, dg *discordgo.Session) ([]map[string]interface{}, error) {
	// defer dg.Close() // fermer avant de return

	list, err := dg.GuildChannels(common.GuildID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de récuperé les fichiers : " + err.Error()})
		return nil, fmt.Errorf("impossible de récupérer les fichiers : %w", err)
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

	return fileChannels, nil
}

func ContainChannel(fileChannels []map[string]interface{}, target string) bool {

	sanitizeTarget := strings.ReplaceAll(strings.ToLower(target), ".", "")

	for _, elem := range fileChannels {
		// fmt.Printf("elem['file_name'] => %s\n", elem["file_name"])
		// fmt.Printf("target => %s\n", target)
		if channel, ok := elem["file_name"].(string); ok && channel == sanitizeTarget {
			return true
		}
	}
	return false
}

func GetChannel(fileChannels []map[string]interface{}, target string) map[string]interface{} {

	sanitizeTarget := strings.ReplaceAll(target, ".", "")

	for _, elem := range fileChannels {
		// fmt.Printf("elem['file_name'] => %s\n", elem["file_name"])
		// fmt.Printf("target => %s\n", target)
		if channel, ok := elem["file_name"].(string); ok && channel == sanitizeTarget {
			return elem
		}
	}
	return nil
}
