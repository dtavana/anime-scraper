package commands

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "add-anime",
			Description: "Add an anime to watchlist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:     "url",
					Type:     discordgo.ApplicationCommandOptionString,
					Required: true,
				},
			},
		},
	}
)
