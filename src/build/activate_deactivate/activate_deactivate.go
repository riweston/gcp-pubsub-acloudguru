package activate_deactivate

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/slack-go/slack"
	"os"
)

// Message structs
type Request struct {
	requestType  string
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
	functions.CloudEvent("ActivateDeactivate", activateDeactivate)
}

func activateDeactivate(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	request := Request{
		requestType:  msg.Message.Attributes["request_type"],
		emailAddress: msg.Message.Attributes["user_email"],
	}

	responseUrl := msg.Message.Attributes["response_url"]
	client := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	slack.MsgOptionReplaceOriginal(responseUrl)
	slackMsg := fmt.Sprintln("User", request.emailAddress, "has been", request.requestType, "by", "admin")
	client.PostMessage("", slack.MsgOptionReplaceOriginal(responseUrl), slack.MsgOptionText(slackMsg, false))

	return nil
}
