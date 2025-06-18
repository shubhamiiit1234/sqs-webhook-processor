// producer.go
package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	messageBody1 := `{
						"webhook_id":"12345",
						"event":"payment_success",
						"data":{
								"amount":250
								}
					}`

	if err := produceMsg(messageBody1); err != nil {
		os.Exit(1)
	}

	messageBody2 := `{
						"webhook_id":"12233",
						"event":"payment_success",
						"data":{
								"amount":250
								}
					}`

	if err := produceMsg(messageBody2); err != nil {
		os.Exit(1)
	}

	messageBody3 := `{
						"webhook_id":"11111",
						"event":"payment_success",
						"data":{
								"amount":250
								}
					}`

	if err := produceMsg(messageBody3); err != nil {
		os.Exit(1)
	}
}

func produceMsg(messageBody string) error {

	queueURL := "https://sqs.ap-south-1.amazonaws.com/400095111419/Queue-Webhook" // Replace with your real queue URL

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"), // Replace with your region
	}))

	sqsClient := sqs.New(sess)

	result, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})

	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}

	fmt.Println("Message sent successfully, ID:", *result.MessageId)
	return nil

}
