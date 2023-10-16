package handlers

import (
	"github.com/dtavana/anime-scraper/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	db           *DatabaseHandler
	notification *NotificationHandler
	dis          *discordgo.Session
}

type CommandHandlerFunctionType = func(s *discordgo.Session, i *discordgo.InteractionCreate)

var (
	integerOptionMinValue = 1.0

	commandData = []*discordgo.ApplicationCommand{
		{
			Name:        "add-anime",
			Description: "Add an anime to watchlist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "URL to add to watchlist",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "episode",
					Description: "Current episode the anime is at",
					MinValue:    &integerOptionMinValue,
					Required:    true,
				},
			},
		},
	}
)

func MakeCommandHandler(db *DatabaseHandler, notification *NotificationHandler, dis *discordgo.Session) *CommandHandler {
	commandHandler := &CommandHandler{db, notification, dis}
	commandHandler.registerHandlers()
	commandHandler.initialize()
	return commandHandler
}

func (c CommandHandler) generateHandlers() map[string]CommandHandlerFunctionType {
	return map[string]CommandHandlerFunctionType{
		"add-anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						util.SuccessEmbed(
							fmt.Sprintf("[Succesfully started tracking new anime starting at episode #%d](%s)", options[1].Value, options[0].Value),
						),
					},
				},
			})
		},
	}
}

func (c CommandHandler) registerHandlers() {
	commandHandlers := c.generateHandlers()
	c.dis.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (c CommandHandler) initialize() {
	c.dis.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := c.dis.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandData))
	for i, v := range commandData {
		cmd, err := c.dis.ApplicationCommandCreate(c.dis.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}
