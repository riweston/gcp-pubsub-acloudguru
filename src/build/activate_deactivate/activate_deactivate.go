package activate_deactivate

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

	// Unmarshal the message
	request := util.UnmarshalActivate(msg)

	/*	clientAcg, err := acg.NewClient(&serviceConfig.ApiKey, &serviceConfig.ConsumerId)
		if err != nil {
			return fmt.Errorf("(ACG Client) Error creating client: %s", err)
		}*/
	userAcg := acg.User{UserId: request.UserId}
	var slackMsg string
	if request.RequestType == "activate" {
		//clientAcg.SetUserActivated(&(userAcg), true)
		slackMsg = fmt.Sprintln("User", userAcg.UserId, "has been activated by admin")

	} else if request.RequestType == "deactivate" {
		//clientAcg.SetUserActivated(&(userAcg), false)
		slackMsg = fmt.Sprintln("User", userAcg.UserId, "has been deactivated by admin")
	}

	clientSlack := slack.New(serviceConfig.SlackBotToken)
	clientSlack.PostMessage("", slack.MsgOptionReplaceOriginal(request.ResponseUrl), slack.MsgOptionText(slackMsg, false))

	return nil
}
