package util

import (
	"github.com/bwmarrin/discordgo"
)

func SuccessEmbed(message string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Description: message,
		Color:       0x00FF00,
	}
	return &embed
}
