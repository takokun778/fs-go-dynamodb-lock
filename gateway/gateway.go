package gateway

import (
	"context"
	"lock/db"
	"lock/model"
	"lock/repository"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var version sync.Map

func init() {
	version = sync.Map{}
}

type Gateway struct {
	database *db.DB
}

func NewGateway(database *db.DB) repository.Repository {
	return &Gateway{
		database: database,
	}
}

type Model struct {
	ID        string `dynamodbav:"id"`
	Name      string `dynamodbav:"name"`
	Version   int    `dynamodbav:"version"`
	CreatedAt string `dynamodbav:"createdAt"`
	UpdatedAt string `dynamodbav:"updatedAt"`
}

func (g *Gateway) Save(ctx context.Context, m model.Model) error {
	now := time.Now().String()

	item, err := attributevalue.MarshalMap(Model{
		ID:        m.ID,
		Name:      m.Name,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("models"),
		Item:      item,
	}

	if _, err := g.database.Client.PutItem(ctx, input); err != nil {
		return err
	}

	return nil
}

func (g *Gateway) Update(ctx context.Context, m model.Model) error {
	rid := GetRIDCtx(ctx)

	v, ok := version.Load(rid)

	if !ok {
		log.Println(ok)
		v = 1
	}

	v = v.(int)

	version := v.(int)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("models"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: m.ID,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#name":      "name",
			"#updatedAt": "updatedAt",
			"#version":   "version",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				Value: m.Name,
			},
			":updatedAt": &types.AttributeValueMemberS{
				Value: time.Now().String(),
			},
			":beforeVersion": &types.AttributeValueMemberN{
				Value: strconv.Itoa(version),
			},
			":afterVersion": &types.AttributeValueMemberN{
				Value: strconv.Itoa(version + 1),
			},
		},
		UpdateExpression:    aws.String("set #name = :name, #updatedAt = :updatedAt, #version = :afterVersion"),
		ConditionExpression: aws.String("#version=:beforeVersion"),
	}

	if _, err := g.database.Client.UpdateItem(ctx, input); err != nil {
		return err
	}

	return nil
}

func (g *Gateway) Find(ctx context.Context, id string) (model.Model, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("models"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	output, err := g.database.Client.GetItem(ctx, input)
	if err != nil {
		return model.Model{}, err
	}

	var m Model

	if err := attributevalue.UnmarshalMap(output.Item, &m); err != nil {
		return model.Model{}, err
	}

	rid := GetRIDCtx(ctx)

	version.Store(rid, m.Version)

	log.Printf("gateway found %+v\n", m)

	return model.Model{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

func (g *Gateway) Delete(ctx context.Context, id string) error {
	return nil
}
