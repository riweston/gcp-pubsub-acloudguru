package pub

import (
	"context"
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type SlashCommand slack.SlashCommand

func (s *SlashCommand) PublishMessage() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic(os.Getenv("TOPIC_NAME"))
	defer topic.Stop()
	r := topic.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"user_id":      s.UserID,
			"response_url": s.ResponseURL,
		},
	})
	s.logPublishMessage(ctx, r)
}

func (s *SlashCommand) logPublishMessage(ctx context.Context, r *pubsub.PublishResult) {
	var results []*pubsub.PublishResult
	results = append(results, r)
	for _, r := range results {
		id, err := r.Get(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Published a message with a message ID: %s\n", id)
	}
}
