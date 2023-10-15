package commands

import "github.com/bwmarrin/discordgo"

type CommandHandler interface  {
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)  
}