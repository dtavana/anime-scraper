package commands

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	db *dynamodb.DynamoDB
}

type ICommandHandler interface {
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}
