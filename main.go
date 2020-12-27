package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
)

func main() {
	// Getting env variables
	var envs map[string]string
	envs, errEnvVars := godotenv.Read(".env")

	if errEnvVars != nil {
		log.Fatal("Error loading .env file")
	}

	// Creating AWS session
	sessionaws, sessionawsErr := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: envs["AWS_PROFILE"],

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(envs["AWS_REGION"]),
		},
	})
	if sessionawsErr != nil {
		fmt.Println("Error to init session")
	}

	fmt.Println("creating AWS session")
	dynamodbsess := dynamodb.New(sessionaws)
	fmt.Println("Listing DynamoDB backups")

	params := dynamodb.ListBackupsInput{}
	//	ctx := context.Background()

	req, resp := dynamodbsess.ListBackupsRequest(&params)
	//req.SetContext(ctx)
	err := req.Send()
	if err == nil { // resp is now filled
		fmt.Println(resp)
	}
}
