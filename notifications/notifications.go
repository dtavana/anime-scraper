package notifications

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func sendNotification(s *discordgo.Session) {
	// Dummy data
	title := "JJK"
	episodeNumber := "690"
	imgUrl := "https://static.bunnycdn.ru/i/cache/images/c/c2/c2c8b3ae50a1b5e71d792ce9cff52431.jpg"
	watchUrl := "https://aniwave.to/watch/jujutsu-kaisen-2nd-season.ll3x3/ep-11"

	s.WebhookExecute(os.Getenv("NOTIFICATION_ID"), os.Getenv("NOTIFICATION_TOKEN"), false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
			Title: "New Episode Detected",
			Description: fmt.Sprintf("%s Episode %s has just released!", title, episodeNumber),
			Color: 0x00FFFF,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: imgUrl,
			},
			URL: watchUrl,
		},
		},
	})
}