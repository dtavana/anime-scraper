package commands

import (
	"github.com/bwmarrin/discordgo" 
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type CommandHandler struct {
	db *dynamodb.DynamoDB
}

type ICommandHandler interface  {
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)  
}