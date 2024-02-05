package webhook_handler

import (
	"context"
	"encoding/json"
	"github.com/alexflint/go-arg"
	"github.com/cloudreach/gcp-pubsub-acloudguru/src/pkg/util"
	"github.com/slack-go/slack"
	"io"
	"net/http"
)

type args struct {
	ProjectID          string `arg:"required,env:PROJECT_ID"`
	TopicName          string `arg:"required,env:TOPIC_NAME"`
	SlackSigningSecret string `arg:"required,env:SLACK_SIGNING_SECRET"`
}

var serviceConfig args

func init() {
	arg.MustParse(&serviceConfig)
}

func EntryPoint(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, serviceConfig.SlackSigningSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := &slack.Msg{Text: "Request received!"}
	b, err := json.Marshal(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		return
	}

	ctx := context.Background()
	clientPubSub, err := util.NewClientPubSub(ctx, serviceConfig.ProjectID, serviceConfig.TopicName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer clientPubSub.Close()

	request := util.NewRequest(s)
	clientPubSub.PublishMessage(ctx, request.NewMessage())
}
