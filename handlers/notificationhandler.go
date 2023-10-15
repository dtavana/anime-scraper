package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type NotificationHandler struct {
	db  *DatabaseHandler
	dis *discordgo.Session
}

func MakeNotificationHandler(db *DatabaseHandler, dis *discordgo.Session) *NotificationHandler {
	return &NotificationHandler{db, dis}
}

func (n NotificationHandler) QueryForAnime(url string) {
	anime := n.db.QueryForAnime(url)
	if anime != nil {
		currentEpisode := 124
		if anime.LastEpisode < currentEpisode {
			n.sendEpisodeNotification(currentEpisode)
		}
		log.Printf("%s %d", anime.Url, anime.LastEpisode)
	}

}

func (n NotificationHandler) sendWebhook(embed *discordgo.MessageEmbed) {
	n.dis.WebhookExecute(os.Getenv("NOTIFICATION_ID"), os.Getenv("NOTIFICATION_TOKEN"), false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}

func (n NotificationHandler) sendEpisodeNotification(currentEpisode int) {
	// Dummy data
	title := "JJK"
	imgUrl := "https://static.bunnycdn.ru/i/cache/images/c/c2/c2c8b3ae50a1b5e71d792ce9cff52431.jpg"
	watchUrl := "https://aniwave.to/watch/jujutsu-kaisen-2nd-season.ll3x3/ep-11"

	embed := discordgo.MessageEmbed{
		Title:       "New Episode Detected",
		Description: fmt.Sprintf("%s Episode #%d has just released!", title, currentEpisode),
		Color:       0x00FFFF,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: imgUrl,
		},
		URL: watchUrl,
	}
	n.sendWebhook(&embed)
}
