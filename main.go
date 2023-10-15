package main

import (
	"animescraper/handlers"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
}

var s *discordgo.Session

func init() {
	var err error
	s, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot instance: %v", err)
	}
}

var (
	db           *handlers.DatabaseHandler
	notification *handlers.NotificationHandler
)

func init() {
	db = handlers.MakeDatabaseHandler()
	notification = handlers.MakeNotificationHandler(db, s)
}

func main() {
	notification.QueryForAnime("test")
}
