package commands

import "github.com/bwmarrin/discordgo"

type AddAnime struct{}

func (a AddAnime) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	})
}