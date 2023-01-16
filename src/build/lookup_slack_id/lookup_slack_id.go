package lookup_slack_id

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	acg "github.com/riweston/acloudguru-client-go"
	"github.com/slack-go/slack"
	"log"
	"os"
)

var apiKey string
var consumerId string
var slackBotToken string

// init functions fail fast if the environment variables are not set

func init() {
	apiKey = os.Getenv("ACLOUDGURU_API_KEY")
	if apiKey == "" {
		fmt.Println("ACLOUDGURU_API_KEY not set")
		os.Exit(1)
	}
	consumerId = os.Getenv("ACLOUDGURU_CONSUMER_ID")
	if consumerId == "" {
		fmt.Println("ACLOUDGURU_CONSUMER_ID not set")
		os.Exit(1)
	}
	slackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	if slackBotToken == "" {
		fmt.Println("SLACK_BOT_TOKEN not set")
		os.Exit(1)
	}
}

// Message structs

type PubSubMsg struct {
	userId      string
	responseUrl string
	requestType string
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

func init() {
	functions.CloudEvent("EntryPoint", entryPoint)
}

func entryPoint(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	// Get the user's email address from Slack
	fmt.Println("Looking up email address for user", msg.Message.Attributes["user_id"])
	slackId, err := lookupEmail(msg.Message.Attributes["user_id"])
	if err != nil {
		return fmt.Errorf("(Slack Client) Error looking up email: %s", err)
	}
	if slackId == "" {
		return fmt.Errorf("(Slack Client) userId is blank")
	}
	fmt.Println("Email address for user", msg.Message.Attributes["user_id"], "is", slackId)

	// Request the user's ACG ID from ACG using email address
	clientAcg, err := acg.NewClient(&apiKey, &consumerId)
	if err != nil {
		return fmt.Errorf("(ACG Client) Error creating client: %s", err)
	}
	acgId, err := clientAcg.GetUserFromEmail(slackId)
	if err != nil {
		return fmt.Errorf("(ACG Client) Error getting user: %s", err)
	}
	println("ACG ID for user", slackId, "is", acgId)
	pubSubMsg := PubSubMsg{
		userId:      (*acgId)[0].UserId,
		responseUrl: msg.Message.Attributes["response_url"],
		requestType: "activate",
	}
	pubSubMsg.PublishMessage()

	return nil
}

func lookupEmail(userId string) (string, error) {
	client := slack.New(slackBotToken)
	user, err := client.GetUserInfo(userId)
	if err != nil {
		log.Printf("Error getting user info: %s", err)
		return "", err
	}
	return user.Profile.Email, nil
}

func (p *PubSubMsg) PublishMessage() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic(os.Getenv("TOPIC_NAME"))
	defer topic.Stop()
	topic.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"user_id":      p.userId,
			"response_url": p.responseUrl,
			"request_type": p.requestType,
		},
	})
}
