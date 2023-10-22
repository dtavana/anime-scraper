package handlers

import (
	"fmt"
	"log"

	"github.com/dtavana/anime-scraper/util"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	db            *DatabaseHandler
	notifications *NotificationHandler
	animes        *AnimeHandler
	dis           *discordgo.Session
}

type CommandHandlerFunctionType = func(s *discordgo.Session, i *discordgo.InteractionCreate)

var (
	commandData = []*discordgo.ApplicationCommand{
		{
			Name:        "add-anime",
			Description: "Add an anime to watchlist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mal-url",
					Description: "MyAnimeList URL for the show",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "watch-url",
					Description: "Watch URL for the show",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete-anime",
			Description: "Delete an anime from watchlist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mal-url",
					Description: "MyAnimeList URL for the show",
					Required:    true,
				},
			},
		},
	}
)

func MakeCommandHandler(db *DatabaseHandler, notification *NotificationHandler, anime *AnimeHandler, dis *discordgo.Session) *CommandHandler {
	commandHandler := &CommandHandler{db, notification, anime, dis}
	commandHandler.registerHandlers()
	commandHandler.initialize()
	return commandHandler
}

func (c CommandHandler) generateHandlers() map[string]CommandHandlerFunctionType {
	return map[string]CommandHandlerFunctionType{
		"add-anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			malUrl, watchUrl := options[0].StringValue(), options[1].StringValue()
			malId, err := util.ParseMalIdFromUrl(malUrl)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							util.ErrorEmbed(
								"Failed to parse ID from supplied MyAnimeList URL",
							),
						},
					},
				})
				return
			}
			animeData := c.animes.QueryAnimeData(malId)
			if animeData == nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							util.ErrorEmbed(
								fmt.Sprintf("Failed to query MyAnimeList data with [the supplied url](%s)", malUrl),
							),
						},
					},
				})
				return
			}
			if c.db.AddAnime(malId, watchUrl) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								URL:         animeData.Data.Url,
								Title:       "New anime added to watchlist",
								Description: fmt.Sprintf("*%s* has been added to the watchlist", animeData.Data.TitleEnglish),
								Color:       util.ColorOrange,
								Image: &discordgo.MessageEmbedImage{
									URL: animeData.Data.Images.Jpg.ImageUrl,
								},
							},
						},
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							util.ErrorEmbed(
								"Failed to save new anime (is it already added to watchlist?)",
							),
						},
					},
				})
			}
		},
		"delete-anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			malUrl := options[0].StringValue()
			malId, err := util.ParseMalIdFromUrl(malUrl)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							util.ErrorEmbed(
								"Failed to parse ID from supplied MyAnimeList URL",
							),
						},
					},
				})
				return
			}
			if c.db.DeleteAnime(malId) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							util.SuccessEmbed(
								fmt.Sprintf("[Succesfully deleted anime from watchlist](%s)", malUrl),
							),
						},
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							util.ErrorEmbed(
								"Failed to deleted anime from watchlist (does it exist in watchlist?)",
							),
						},
					},
				})
			}
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

	log.Println("Updating status...")
	c.dis.UpdateWatchStatus(0, "your favorite animes!")

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
