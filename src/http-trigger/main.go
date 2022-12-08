package hello_world

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
)

// HelloWorld prints "Hello, world."
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	fmt.Fprintln(w, "Hello, world.")
	client, err := pubsub.NewClient(ctx, "cr-lab-rweston-2811223416")
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic("acg-request")
	topic.Publish(ctx, &pubsub.Message{
		Data: []byte("hello world"),
		Attributes: map[string]string{
			"user":     "me!",
			"activate": "false",
		},
	})
}
