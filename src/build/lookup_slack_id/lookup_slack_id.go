package lookup_slack_id

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/alexflint/go-arg"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudreach/gcp-pubsub-acloudguru/src/pkg/util"
	acg "github.com/riweston/acloudguru-client-go"
	"github.com/slack-go/slack"
)

// Required environment variables

type args struct {
	ProjectID     string `arg:"required,env:PROJECT_ID"`
	TopicName     string `arg:"required,env:TOPIC_NAME"`
	ApiKey        string `arg:"required,env:ACLOUDGURU_API_KEY"`
	ConsumerId    string `arg:"required,env:ACLOUDGURU_CONSUMER_ID"`
	SlackBotToken string `arg:"required,env:SLACK_BOT_TOKEN"`
}

var serviceConfig args

func init() {
	arg.MustParse(&serviceConfig)
}

func init() {
	functions.CloudEvent("EntryPoint", entryPoint)
}

func entryPoint(ctx context.Context, e event.Event) error {
	var msg util.MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	// Get the user's email address from Slack
	clientSlack := slack.New(serviceConfig.SlackBotToken)
	slackId, err := clientSlack.GetUserInfo(msg.Message.Attributes["user_id"])
	if err != nil {
		return fmt.Errorf("(Slack Client) Error getting user: %s", err)
	}

	// Request the user's ACG ID from ACG using email address
	clientAcg, err := acg.NewClient(&serviceConfig.ApiKey, &serviceConfig.ConsumerId)
	if err != nil {
		return fmt.Errorf("(ACG Client) Error creating client: %s", err)
	}
	acgId, err := clientAcg.GetUserFromEmail(slackId.Profile.Email)
	if err != nil {
		return fmt.Errorf("(ACG Client) Error getting user: %s", err)
	}

	// Publish the message to PubSub
	clientPubSub, err := util.NewClientPubSub(ctx, serviceConfig.ProjectID, serviceConfig.TopicName)
	if err != nil {
		fmt.Errorf("error creating client: %v", err)
	}
	defer clientPubSub.Close()
	pubSubMsg := util.NewActivate((*acgId)[0], msg)
	clientPubSub.PublishMessage(ctx, pubSubMsg.NewMessage())

	return nil
}
