package upload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fichier manquant"})
		return
	}

	tempFile, err := os.Create(filepath.Join(os.TempDir(), file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur l'ors de la création du fichier"})
		return
	}
	defer tempFile.Close()

	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur l'ors de l'enregistrement du fichier"})
		return
	}

	channelID, err := createChannelAndSegments(tempFile.Name(), file.Filename, *file, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du canal ou de l'envoi des segments : " + err.Error()})
		return
	}

	os.Remove(tempFile.Name())

	c.JSON(http.StatusOK, gin.H{"message": "Fichier uploadé et segments envoyés", "channel_id": channelID})

}

func createChannelAndSegments(filePath, filename string, file multipart.FileHeader, c *gin.Context) (string, error) {
	dg, err := discordgo.New("Bot " + common.Token)
	guildID := c.Param("guildID")

	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du client Discord: %v", err)
	}

	err = dg.Open()
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'ouverture de la connexion: %v", err)
	}
	defer dg.Close()

	channelList, err := list.ListChannelFileWithDg(c, dg)

	if err != nil {
		return "", fmt.Errorf("erreur l'ors de la récuperations des channels: %v", err)
	}

	date := time.Now().Format("20060102-150405")

	ext := filepath.Ext(filename)
	filename = filename[:len(filename)-len(ext)]

	channelName := fmt.Sprintf("%s__%s__%d__%s__%s", uuid.New().String(), filename, file.Size, date, ext)

	if list.ContainChannel(channelList, filename) {
		return "", fmt.Errorf("erreur el fichier existe déja: %v", err)
	}

	channel, err := dg.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du canal: %v", err)
	}

	if err := sendFileSegmentToChannel(filePath, channel.ID, dg); err != nil {
		return "", err
	}

	return channel.ID, nil
}

func sendFileSegmentToChannel(filePath, channelID string, dg *discordgo.Session) error {
	file, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 10*1024*1024) // 10 Mo
	segmentIndex := 0

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erreur lors de la lecture du fichier: %v", err)
		}

		fileSegment := bytes.NewReader(buffer[:n])
		_, err = dg.ChannelFileSend(channelID, fmt.Sprintf("segment_%d.dat", segmentIndex), fileSegment)
		if err != nil {
			return fmt.Errorf("erreur lors de l'envoi du segment %d: %v", segmentIndex, err)
		}

		segmentIndex++
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
