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
	queueURL := "https://sqs.ap-south-1.amazonaws.com/400095111419/Queue-Webhook" // Replace with your real queue URL

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"), // Replace with your region
	}))

	sqsClient := sqs.New(sess)

	messageBody := `{
						"webhook_id":"12345",
						"event":"payment_success",
						"data":{
								"amount":250
								}
					}`

	result, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})

	if err != nil {
		fmt.Println("Error sending message:", err)
		os.Exit(1)
	}

	fmt.Println("Message sent successfully, ID:", *result.MessageId)
}
