package upload

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	Token   = os.Getenv("DISCORD_TOKEN")
	GuildID = "1340725915130400859"
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

	channelID, err := createChannelAndSegments(tempFile.Name(), file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du canal ou de l'envoi des segments : " + err.Error()})
		return
	}

	os.Remove(tempFile.Name())

	c.JSON(http.StatusOK, gin.H{"message": "Fichier uploadé et segments envoyés", "channel_id": channelID})

}

func createChannelAndSegments(filePath, filename string) (string, error) {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du client Discord: %v", err)
	}

	err = dg.Open()
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'ouverture de la connexion: %v", err)
	}
	defer dg.Close()

	channel, err := dg.GuildChannelCreate(GuildID, uuid.New().String()+"__"+filename, discordgo.ChannelTypeGuildText)
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
