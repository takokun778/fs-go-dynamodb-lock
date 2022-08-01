package ddl

import (
	"context"
	"lock/db"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func init() {
	db, err := db.NewDB(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	ddl := NewDDL(db)

	if err := ddl.CreateTableModel(context.Background()); err != nil {
		log.Println(err.Error())
	}
}

type DDL struct {
	database *db.DB
}

func NewDDL(database *db.DB) *DDL {
	return &DDL{
		database: database,
	}
}

func (d *DDL) CreateTableModel(ctx context.Context) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("models"),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}

	if _, err := d.database.Client.CreateTable(ctx, input); err != nil {
		return err
	}

	return nil
}

func (d *DDL) DeleteTableModel(ctx context.Context) error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("models"),
	}

	if _, err := d.database.Client.DeleteTable(ctx, input); err != nil {
		return err
	}

	return nil
}

func (d *DDL) ListTables(ctx context.Context) ([]string, error) {
	input := &dynamodb.ListTablesInput{}

	output, err := d.database.Client.ListTables(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.TableNames, nil
}
