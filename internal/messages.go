package internal

import (
	"cloud.google.com/go/pubsub"
	"context"
	acg "github.com/riweston/acloudguru-client-go"
	"github.com/slack-go/slack"
)

type ClientPubSub struct {
	Client    *pubsub.Client
	TopicName string
}

type Request struct {
	UserId      string
	ResponseUrl string
}

type Activate struct {
	UserId      string
	ResponseUrl string
	RequestType string
}

func NewClientPubSub(ctx context.Context, projectId string, TopicName string) (*ClientPubSub, error) {
	clientPubSub, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return &ClientPubSub{
		Client:    clientPubSub,
		TopicName: TopicName,
	}, nil
}

func (c *ClientPubSub) Close() error {
	return c.Client.Close()
}

func NewDeactivate(user acg.User, data MessagePublishedData) *Activate {
	return &Activate{
		UserId:      user.UserId,
		ResponseUrl: data.Message.Attributes["response_url"],
		RequestType: "deactivate",
	}
}

func NewActivate(user acg.User, data MessagePublishedData) *Activate {
	return &Activate{
		UserId:      user.UserId,
		ResponseUrl: data.Message.Attributes["response_url"],
		RequestType: "activate",
	}
}

func (r *Activate) NewMessage() *pubsub.Message {
	return &pubsub.Message{
		Attributes: map[string]string{
			"user_id":      r.UserId,
			"response_url": r.ResponseUrl,
			"request_type": r.RequestType,
		},
	}
}

func UnmarshalActivate(data MessagePublishedData) *Activate {
	return &Activate{
		UserId:      data.Message.Attributes["user_id"],
		RequestType: data.Message.Attributes["request_type"],
		ResponseUrl: data.Message.Attributes["response_url"],
	}
}

func NewRequest(s slack.SlashCommand) *Request {
	return &Request{
		UserId:      s.UserID,
		ResponseUrl: s.ResponseURL,
	}
}

func (r *Request) NewMessage() *pubsub.Message {
	return &pubsub.Message{
		Attributes: map[string]string{
			"user_id":      r.UserId,
			"response_url": r.ResponseUrl,
		},
	}
}

func (c *ClientPubSub) PublishMessage(ctx context.Context, m *pubsub.Message) {
	topic := c.Client.Topic(c.TopicName)
	defer topic.Stop()
	topic.Publish(ctx, m)
}
