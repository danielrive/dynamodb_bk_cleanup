package main

import (
	"fmt"
	"log"
	"time"

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
	for {
		if listDynamoBackups(dynamodbsess) != 0 {
			removeDynamoBackup(dynamodbsess)
		} else {
			break
		}
	}
}

func removeDynamoBackup(awsession *dynamodb.DynamoDB) {
	inputlist := dynamodb.ListBackupsInput{
		Limit: aws.Int64(10),
	}
	req, resp := awsession.ListBackupsRequest(&inputlist)
	err := req.Send()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(resp.BackupSummaries); i++ {
		reqdelete, respdelete := awsession.DeleteBackupRequest(&dynamodb.DeleteBackupInput{
			BackupArn: resp.BackupSummaries[i].BackupArn,
		})

		err := reqdelete.Send()
		if err != nil { // resp is now filled
			panic(err)
		}
		fmt.Println(respdelete)
		//	arnBackups = append(arnBackups, *resp.BackupSummaries[i].BackupArn)
	}
	//fmt.Println(arnBackups)
	//fmt.Println(len(arnBackups))

}

func listDynamoBackups(awsession *dynamodb.DynamoDB) int {
	// code from https://stackoverflow.com/questions/54118604/dynamodb-list-all-backups-using-aws-golang-sdk
	var exclusiveStartARN *string
	var backups []*dynamodb.BackupSummary
	fmt.Println("Listing DynamoDB backups")

	for {
		backup, err := awsession.ListBackups(&dynamodb.ListBackupsInput{
			ExclusiveStartBackupArn: exclusiveStartARN,
		})
		if err != nil {
			panic(err)
		}
		backups = append(backups, backup.BackupSummaries...)
		if backup.LastEvaluatedBackupArn != nil {
			exclusiveStartARN = backup.LastEvaluatedBackupArn

			//max 5 times a second so we dont hit the limit
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}
	return len(backups)
}

func createDynamoBackup(awsession *dynamodb.DynamoDB, bkname string, tableName string) {
	params := dynamodb.CreateBackupInput{
		BackupName: aws.String(bkname),
		TableName:  aws.String(tableName),
	}
	req, resp := awsession.CreateBackupRequest(&params)

	err := req.Send()
	if err == nil { // resp is now filled
		fmt.Println(resp)
	}
}
