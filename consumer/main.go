// consumer.go
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Webhook struct {
	WebhookID string `json:"webhook_id"`
	Event     string `json:"event"`
}

func processWebhook(msg string, counter int) error {
	var webhook Webhook
	if err := json.Unmarshal([]byte(msg), &webhook); err != nil {
		return err
	}

	fmt.Println("Processing webhook:", webhook)

	// Simulate failure randomly
	// if rand.Intn(3) == 0 {
	// 	return fmt.Errorf("simulated failure")
	// }

	if counter < 2 {
		counter++
		return fmt.Errorf("simulated failure")
	}

	fmt.Println("Processed webhook successfully!")
	return nil
}

func main() {
	queueURL := "https://sqs.ap-south-1.amazonaws.com/400095111419/Queue-Webhook" // Replace with your real queue URL

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	}))

	sqsClient := sqs.New(sess)
	rand.Seed(time.Now().UnixNano())

	counter := 0

	for {
		output, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(10),
		})

		if err != nil {
			fmt.Println("Receive error:", err)
			continue
		}

		if len(output.Messages) == 0 {
			fmt.Println("No messages.")
			continue
		}

		for _, msg := range output.Messages {
			fmt.Println("Received message:", *msg.Body)

			if err := processWebhook(*msg.Body, counter); err != nil {
				fmt.Println("Failed to process:", err)
				// Don’t delete — it will retry and go to DLQ eventually
				continue
			}

			// Delete if processed successfully
			_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				fmt.Println("Error deleting message:", err)
			} else {
				fmt.Println("Deleted message successfully.")
			}
		}
	}
}
