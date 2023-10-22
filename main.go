package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dtavana/anime-scraper/handlers"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
}

var dis *discordgo.Session

func init() {
	var err error
	dis, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot instance: %v", err)
	}
}

var (
	db           *handlers.DatabaseHandler
	notification *handlers.NotificationHandler
	animes       *handlers.AnimeHandler
	command      *handlers.CommandHandler
)

func init() {
	db = handlers.MakeDatabaseHandler()
	notification = handlers.MakeNotificationHandler(db, dis)
	animes = handlers.MakeAnimeHandler()
	command = handlers.MakeCommandHandler(db, notification, animes, dis)
}

func main() {
	defer dis.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
