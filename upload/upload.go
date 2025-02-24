package upload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"discord_drive/common"
	"discord_drive/list"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	guildID := c.Param("guildID")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fichier manquant"})
		return
	}

	fmt.Printf("Taille totale du fichier {%s}: %.2f Mb\n", file.Filename, float64(file.Size)/1024/1024)
	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ouverture du fichier : " + err.Error()})
		return
	}
	// defer srcFile.Close()

	c.JSON(http.StatusAccepted, gin.H{"message": "Fichier en cours de traitement"})
	go func() {
		channelList, err := list.ListChannelFileWithDg(c, common.DiscordSession)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récuperation des channels : " + err.Error()})
			return
		}
		defer srcFile.Close() // fermer le file a la fin de la routine
		channelID, err := createChannelAndSegments(srcFile, file.Filename, *file, guildID, channelList)
		if err != nil {
			fmt.Printf("❌ Erreur lors de la création du canal ou de l'envoi des segments : %v\n", err)
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du canal ou de l'envoi des segments : " + err.Error()})
			return
		}
		fmt.Printf("✅ Fichier traité, channel ID: %s\n", channelID)
	}()

	// c.JSON(http.StatusOK, gin.H{"message": "Fichier uploadé et segments envoyés", "channel_id": channelID})
}

func createChannelAndSegments(reader io.Reader, filename string, file multipart.FileHeader, guildID string, channelList []map[string]interface{}) (string, error) {
	dg := common.DiscordSession

	date := time.Now().Format("20060102-150405")

	ext := filepath.Ext(filename)
	filename = filename[:len(filename)-len(ext)]

	channelName := fmt.Sprintf("%s__%s__%d__%s__%s", uuid.New().String(), filename, file.Size, date, ext)

	if list.ContainChannel(channelList, filename) {
		return "", fmt.Errorf("erreur el fichier existe déja")
	}

	channel, err := dg.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du canal: %v", err)
	}

	if err := sendFileSegmentToChannel(reader, channel.ID, dg); err != nil {
		return "", err
	}

	return channel.ID, nil
}

func sendFileSegmentToChannel(reader io.Reader, channelID string, dg *discordgo.Session) error {
	buffer := make([]byte, 10*1024*1024) // 10 Mo
	segmentIndex := 0

	for {
		// lire depuis le stream
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("erreur lors de la lecture du fichier: %v", err)
		}

		fileSegment := bytes.NewReader(buffer[:n])
		fmt.Printf("Tentative d'envoi du segment %d de taille: %d bytes\n", segmentIndex, n)
		_, err = dg.ChannelFileSend(channelID, fmt.Sprintf("segment_%d.dat", segmentIndex), fileSegment)
		if err != nil {
			return fmt.Errorf("erreur lors de l'envoi du segment %d: %v", segmentIndex, err)
		}

		segmentIndex++
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}
