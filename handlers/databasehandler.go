package handlers

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Url string
	LastEpisode int
}

type DatabaseHandler struct {
	session *session.Session
	svc *dynamodb.DynamoDB
}

func MakeDatabaseHandler() *DatabaseHandler {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://localhost:8000"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	svc := dynamodb.New(session)
	handler := &DatabaseHandler{session, svc}
	handler.setupDatabase()
	return handler
}

func (d DatabaseHandler) setupDatabase() {
	d.session = session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://localhost:8000"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	d.svc = dynamodb.New(d.session)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("url"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("url"),
				KeyType: aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(os.Getenv("DB_TABLE")),
	}
	d.svc.CreateTable(input)
}

func (d DatabaseHandler) QueryForAnime(url string) *Item {
	res, err := d.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"url": {
				S: aws.String(url),
			},
		},
		TableName: aws.String(os.Getenv("DB_TABLE")),
	})
	if err != nil {
		log.Printf("Error querying for anime %v", err)
	}
	if res.Item == nil {
		log.Printf("Could not find item with URL: %s", url)
	} else {
		item := Item{}
		err = dynamodbattribute.UnmarshalMap(res.Item, &item)
		if err != nil {
			log.Printf("Error unmarshalling item %v", err)
		}
		return &item
	}
	return nil
}
