package comm

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Text    string
	HTML    string
}

var AWS_REGION = "us-east-1"

func newEmail() *Email {
	return &Email{}
}

func (e *Email) send() (*ses.SendEmailOutput, error) {
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" {
		log.Fatal("AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and TO_EMAIL must be set")
	}

	// To make calls to an AWS service, you must first
	// construct a service client instance. A service client
	// provides low-level access to every API action for that
	// service. For example, you create an Amazon SES service
	// client to make calls to Amazon SES APIs.
	//
	// config.LoadDefaultConfig(context.TODO()) will construct
	// an aws.Config using the AWS shared configuration sources.
	// https://docs.aws.amazon.com/sdkref/latest/guide/overview.html
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := ses.NewFromConfig(cfg)

	email, err := client.SendEmail(context.TODO(), &ses.SendEmailInput{
		Source: aws.String(e.From),
		Destination: &types.Destination{
			ToAddresses: e.To,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(e.Subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(e.Text),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return email, nil
}
