package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
)

type Payload struct {
	Choices           []Choice `json:"choices"`
	Created           int      `json:"created"`
	Id                string   `json:"id"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

var payload *Payload = nil

func main() {

	doPubSub := os.Getenv("PUBSUB_ON_WAKE")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	r := gin.Default()

	r.GET("/sleep", func(c *gin.Context) {
		timeInMs := c.DefaultQuery("ms", "500")
		num, _ := strconv.Atoi(timeInMs)
		time.Sleep(time.Duration(num) * time.Millisecond)

		var id = ""
		if doPubSub == "TRUE" {
			pubId, err := publishMessage(os.Getenv("PROJECT_ID"), os.Getenv("TOPIC_ID"), "WAKE UP")
			id = pubId
			if err != nil {
				fmt.Println(err)
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.String(200, "yawn, slept for "+timeInMs+"ms and published wakeup message "+id)
	})

	r.POST("/payload", func(c *gin.Context) {

		if payload == nil {
			b, err := os.ReadFile("payload.json")
			if err == nil {
				json.Unmarshal(b, payload)
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.JSON(200, *payload)
	})

	r.POST("/payload/:sizeInMb", func(c *gin.Context) {

		sizeInMb := c.Param("sizeInMb")
		sizeNum, err := strconv.Atoi(sizeInMb)
		var largePayload Payload

		if err != nil {
			fmt.Println("Could not convert size parameter to a number - " + sizeInMb + ". Using default 5 megabytes.")
			sizeNum = 5
		} else if sizeNum > 20 {
			fmt.Println("Size parameter " + sizeInMb + " greater than 20mb, using max size 20mb.")
		}

		sizeNum = sizeNum * 350

		b, err := os.ReadFile("payload.json")
		if err == nil {
			json.Unmarshal(b, &largePayload)
		}
		orginalMessage := largePayload.Choices[0].Message.Content
		for sizeNum > 0 {
			largePayload.Choices[0].Message.Content += orginalMessage
			sizeNum--
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.JSON(200, largePayload)
	})

	r.Run(":" + PORT)
}

func publishMessage(projectID, topicID, msg string) (string, error) {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("pubsub: result.Get: %w", err)
	}
	return id, nil
}
