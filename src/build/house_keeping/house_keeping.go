package house_keeping

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/alexflint/go-arg"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudreach/gcp-pubsub-acloudguru/src/pkg/util"
	"github.com/riweston/acloudguru-client-go"
)

// Required environment variables

type args struct {
	ProjectID  string `arg:"required,env:PROJECT_ID"`
	TopicName  string `arg:"required,env:TOPIC_NAME"`
	ApiKey     string `arg:"required,env:ACLOUDGURU_API_KEY"`
	ConsumerId string `arg:"required,env:ACLOUDGURU_CONSUMER_ID"`
	DaysCap    int    `arg:"required,env:DAYS_CAP"`
	LicenseCap int    `arg:"required,env:LICENSE_CAP"`
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

	client, err := acg.NewClient(&serviceConfig.ApiKey, &serviceConfig.ConsumerId)
	if err != nil {
		fmt.Errorf("error creating client: %v", err)
	}
	usersAll, err := client.GetUsersAll()
	if err != nil {
		fmt.Errorf("error getting users: %v", err)
	}
	usersProcessed := util.NewUsers(usersAll, serviceConfig.DaysCap, serviceConfig.LicenseCap)
	usersDeactivate := usersProcessed.GetUsersToDeactivate()

	if len(usersDeactivate) > 0 {
		clientPubSub, err := util.NewClientPubSub(ctx, serviceConfig.ProjectID, serviceConfig.TopicName)
		if err != nil {
			fmt.Errorf("error creating client: %v", err)
		}
		defer clientPubSub.Close()
		for _, user := range usersDeactivate {
			// TODO: Debug statement this should be replaced with proper logging
			fmt.Println(user.Status, user.LastSeenTimestamp, user.Name)
			pubSubMsg := util.NewDeactivate(user, msg)
			fmt.Println(pubSubMsg.UserId, pubSubMsg.RequestType, pubSubMsg.ResponseUrl)
			clientPubSub.PublishMessage(ctx, pubSubMsg.NewMessage())
		}
	}
	return nil
}
