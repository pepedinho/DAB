package list

import (
	"discord_drive/common"
	"fmt"
	"net/http"
	"strconv"
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
	dg := common.DiscordSession
	guildID := c.Param("guildID")
	


	list, err := dg.GuildChannels(guildID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de récuperé les fichiers : " + err.Error()})
		return nil, fmt.Errorf("impossible de récupérer les fichiers : %w", err)
	}

	var fileChannels []map[string]interface{}
	for _, channel := range list {
		if channel.Type == discordgo.ChannelTypeGuildText {
			parts := strings.SplitN(channel.Name, "__", 5)
			if len(parts) >= 4 {
				sizeBytes, err := strconv.Atoi(parts[2])
				if err != nil {
					sizeBytes = 0
				}

				fileChannels = append(fileChannels, map[string]interface{}{
					"id":        channel.ID,
					"channel":   channel.Name,
					"file_name": parts[1],
					"size":      parts[2],
					"mb_size":   fmt.Sprintf("%.1f", float64(sizeBytes)/(1024*1024)),
					"date":      parts[3],
					"extension": parts[4],
				})
			}
		}
	}

	return fileChannels, nil
}

// dg has to be already open before call this function
func ListChannelFileWithDg(c *gin.Context, dg *discordgo.Session) ([]map[string]interface{}, error) {
	// defer dg.Close() // fermer avant de return
	guildID := c.Param("guildID")

	list, err := dg.GuildChannels(guildID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de récuperé les fichiers : " + err.Error()})
		return nil, fmt.Errorf("impossible de récupérer les fichiers : %w", err)
	}

	var fileChannels []map[string]interface{}
	for _, channel := range list {
		if channel.Type == discordgo.ChannelTypeGuildText {
			parts := strings.SplitN(channel.Name, "__", 5)
			if len(parts) >= 4 {
				sizeBytes, err := strconv.Atoi(parts[2])
				if err != nil {
					sizeBytes = 0
				}

				fileChannels = append(fileChannels, map[string]interface{}{
					"id":        channel.ID,
					"channel":   channel.Name,
					"file_name": parts[1],
					"size":      parts[2],
					"mb_size":   fmt.Sprintf("%.1f", float64(sizeBytes)/(1024*1024)),
					"date":      parts[3],
					"extension": parts[4],
				})
			}
		}
	}

	return fileChannels, nil
}

func ContainChannel(fileChannels []map[string]interface{}, target string) bool {

	sanitizeTarget := strings.ReplaceAll(strings.ToLower(target), ".", "")

	// fmt.Printf("Searching => %s | %s\n", sanitizeTarget, target)

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

	sanitizeTarget := strings.ReplaceAll(strings.ToLower(target), ".", "")
	// fmt.Printf("Searching => %s\n", sanitizeTarget)

	for _, elem := range fileChannels {
		// fmt.Printf("elem['file_name'] => %s\n", elem["file_name"])
		// fmt.Printf("target => %s\n", target)
		if channel, ok := elem["file_name"].(string); ok && channel == sanitizeTarget {
			return elem
		}
	}
	return nil
}
