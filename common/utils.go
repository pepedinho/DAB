package common

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token = os.Getenv("DISCORD_TOKEN")
	// GuildID = "1340725915130400859"
	DiscordSession *discordgo.Session
	HttpClient     = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func InitDiscordSession() {
	var err error
	DiscordSession, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client Discord: %v", err)
	}
	DiscordSession.Client = HttpClient

	err = DiscordSession.Open()
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la connexion Discord: %v", err)
	}

	log.Println("✅ Connexion Discord établie !")
}
