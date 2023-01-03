package lookup_slack_id

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/slack-go/slack"
	"log"
	"os"
)

// Message structs
type Request struct {
	userId       string
	responseUrl  string
	emailAddress string
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

func init() {
	functions.CloudEvent("LookupSlackId", lookupSlackId)
}

func lookupSlackId(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	response := Request{
		userId:      msg.Message.Attributes["user_id"],
		responseUrl: msg.Message.Attributes["response_url"],
	}
	if response.userId == "" {
		log.Printf("userId is blank!")
	}
	if err := response.lookupEmail(); err != nil {
		log.Printf("Error looking up email: %s", err)
	}
	log.Printf("Email is %s!", response.emailAddress)
	response.PublishMessage()
	return nil
}

func (r *Request) lookupEmail() error {
	client := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	user, err := client.GetUserInfo(r.userId)
	if err != nil {
		log.Printf("Error getting user info: %s", err)
	}
	r.emailAddress = user.Profile.Email
	return nil
}

func (r *Request) PublishMessage() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic(os.Getenv("TOPIC_NAME"))
	defer topic.Stop()
	topic.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"user_email":   r.emailAddress,
			"response_url": r.responseUrl,
			"request_type": "activate",
		},
	})
}
