package dynamodb

import (
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/belljustin/golinks/pkg/golinks"
)

const (
	linksTableName = "Links"
)

func init() {
	golinks.RegisterStorage("dynamodb", newStorage)
}

type Storage struct {
	tableName string
	svc       *aws_dynamodb.DynamoDB
}

func newStorage() golinks.Storage {
	loadConfig()

	config := &aws.Config{
		Region:   aws.String(C.Storage.Region),
		Endpoint: aws.String(C.Storage.Endpoint),
	}
	sess := session.Must(session.NewSession(config))
	svc := aws_dynamodb.New(sess)

	return &Storage{C.Storage.TableName, svc}
}

func (s Storage) GetLink(name string) (*url.URL, error) {
	linkKey := LinkKey{
		Name: name,
	}

	key, err := dynamodbattribute.MarshalMap(linkKey)
	if err != nil {
		return nil, err
	}

	input := &aws_dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(linksTableName),
	}

	result, err := s.svc.GetItem(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == aws_dynamodb.ErrCodeResourceNotFoundException {
			return nil, nil
		}
		return nil, err
	}

	link := Link{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &link)
	if err != nil {
		return nil, err
	}

	return url.Parse(link.URL)
}

func (s Storage) SetLink(name string, url url.URL) error {
	link := Link{
		Name: name,
		URL:  url.String(),
	}

	info, err := dynamodbattribute.MarshalMap(link)
	if err != nil {
		return err
	}

	input := &aws_dynamodb.PutItemInput{
		Item:      info,
		TableName: aws.String(linksTableName),
	}

	_, err = s.svc.PutItem(input)
	return err
}

func (s Storage) Migrate() error {
	input := &aws_dynamodb.CreateTableInput{
		AttributeDefinitions: []*aws_dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*aws_dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &aws_dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(linksTableName),
	}

	result, err := s.svc.CreateTable(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == aws_dynamodb.ErrCodeResourceInUseException {
			log.Printf("[INFO] dynamodb: table already exists")
			return nil
		}
		panic(err)
	}

	log.Printf("[INFO] dynamodb: completed migration: %s", result)
	return nil
}

type Link struct {
	Name string
	URL  string
}

type LinkKey struct {
	Name string
}
