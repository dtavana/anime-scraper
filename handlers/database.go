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
	MalId   string
	WatchUrl string
}

type DatabaseHandler struct {
	session   *session.Session
	svc       *dynamodb.DynamoDB
	tableName *string
}

func MakeDatabaseHandler() *DatabaseHandler {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://localhost:8000"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(session)
	tableName := aws.String(os.Getenv("DB_TABLE"))
	handler := &DatabaseHandler{session, svc, tableName}
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
				AttributeName: aws.String("MalId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("MalId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: d.tableName,
	}
	d.svc.CreateTable(input)
}

func (d DatabaseHandler) AddAnime(malId string, watchUrl string) bool {
	item := Item{malId, watchUrl}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return false
	}
	_, err = d.svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: d.tableName,
	})
	return err == nil
}

func (d DatabaseHandler) DeleteAnime(malId string) bool {
	_, err := d.svc.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"MalId": {
				S: aws.String(malId),
			},
		},
		TableName: d.tableName,
	})
	return err == nil
}

func (d DatabaseHandler) QueryForAnime(malId string) *Item {
	res, err := d.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"MalId": {
				S: aws.String(malId),
			},
		},
		TableName: d.tableName,
	})
	if err != nil {
		log.Printf("Error querying for anime %v", err)
		return nil
	}
	if res.Item == nil {
		log.Printf("Could not find item with MalID: %s", malId)
		return nil
	} else {
		item := Item{}
		err = dynamodbattribute.UnmarshalMap(res.Item, &item)
		if err != nil {
			log.Printf("Error unmarshalling item %v", err)
		}
		return &item
	}
}
