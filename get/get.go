package get

import (
	"discord_drive/common"
	"discord_drive/list"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	filename := c.Query("filename")

	dg := common.DiscordSession

	channelList, err := list.ListChannelFileWithDg(c, dg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de récupérer les channels : " + err.Error()})
		return
	}

	if !list.ContainChannel(channelList, filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le fichier (" + filename + ") n'existe pas"})
		return
	}

	channel := list.GetChannel(channelList, filename)

	var allMessages []*discordgo.Message

	lastMessageID := ""

	channelID := channel["id"].(string)

	for {
		messages, err := dg.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des messages"})
			return
		}

		if len(messages) == 0 {
			break
		}

		allMessages = append(allMessages, messages...)
		lastMessageID = messages[len(messages)-1].ID
	}

	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].Attachments[0].Filename < allMessages[j].Attachments[0].Filename
	})

	tempFilePath := filepath.Join(os.TempDir(), filename)
	outputFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du fichier temporaire"})
		return
	}
	defer outputFile.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	for _, msg := range allMessages {
		if len(msg.Attachments) == 0 {
			continue
		}
		segmentURL := msg.Attachments[0].URL

		resp, err := http.Get(segmentURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de téléchargement d'un segment"})
			return
		}
		defer resp.Body.Close()
		
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la reconstruction du fichier"})
			return
		}
	}
}
