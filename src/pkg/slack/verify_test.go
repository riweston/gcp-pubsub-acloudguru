package slack

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestVerifyTimeStamp(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Provided with a valid timestamp",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{strconv.FormatInt(time.Now().Unix(), 10)},
					},
				},
			},
			want: true,
		},
		{
			name: "Provided with an old ( > 5 minutes ) timestamp",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{strconv.FormatInt(time.Now().Add(time.Minute*-6).Unix(), 10)},
					},
				},
			},
			want: false,
		},
		{
			name: "Provided with an old timestamp",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{"1670926904.311199"},
					},
				},
			},
			want: false,
		},
		{
			name: "Provided without 'X-Slack-Request-Timestamp' header",
			args: args{
				r: &http.Request{},
			},
			want: false,
		},
		{
			name: "Provided with a 'X-Slack-Request-Timestamp' silly header value",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{"something silly"},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyTimeStamp(tt.args.r); got != tt.want {
				t.Errorf("VerifyTimeStamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	// Taken from https://api.slack.com/authentication/verifying-requests-from-slack
	slackSigningSecret  = "8f742231b10e8888abcd99yyyzzz85a5"
	requestBody         = "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c"
	requestTimestamp    = "1531420618"
	sigBaseString       = "v0:" + requestTimestamp + ":" + requestBody
	sigBaseStringHashed = "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"
)

func TestHashSigBaseString(t *testing.T) {
	type args struct {
		sigBaseString string
		envVarSet     bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Provided with a valid sigBaseString",
			args: args{
				sigBaseString: sigBaseString,
				envVarSet:     true,
			},
			want:    sigBaseStringHashed,
			wantErr: false,
		},
		{
			name: "Environment variable 'SLACK_SIGNING_SECRET' not set",
			args: args{
				sigBaseString: sigBaseString,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envVarSet {
				defer os.Unsetenv("SLACK_SIGNING_SECRET")
				os.Setenv("SLACK_SIGNING_SECRET", slackSigningSecret)
			}
			got, err := HashSigBaseString(tt.args.sigBaseString)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashSigBaseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashSigBaseString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSigBaseString(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Provided with a valid request",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{requestTimestamp},
					},
					Body: io.NopCloser(bytes.NewBufferString(requestBody)),
				},
			},
			want: sigBaseString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSigBaseString(tt.args.r); got != tt.want {
				t.Errorf("GenerateSigBaseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifySignature(t *testing.T) {
	type args struct {
		r              *http.Request
		envVarSet      bool
		envVarMisMatch bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Provided with a valid request",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{requestTimestamp},
						"X-Slack-Signature":         []string{sigBaseStringHashed},
					},
					Body: io.NopCloser(bytes.NewBufferString(requestBody)),
				},
				envVarSet: true,
			},
			wantErr: false,
		},
		{
			name: "Provided without environment variable 'SLACK_SIGNING_SECRET' set",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{requestTimestamp},
						"X-Slack-Signature":         []string{sigBaseStringHashed},
					},
					Body: io.NopCloser(bytes.NewBufferString(requestBody)),
				},
			},
			wantErr: true,
		},
		{
			name: "Provided with incorrect 'SLACK_SIGNING_SECRET' set",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Slack-Request-Timestamp": []string{requestTimestamp},
						"X-Slack-Signature":         []string{sigBaseStringHashed},
					},
					Body: io.NopCloser(bytes.NewBufferString(requestBody)),
				},
				envVarMisMatch: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.envVarSet {
				defer os.Unsetenv("SLACK_SIGNING_SECRET")
				os.Setenv("SLACK_SIGNING_SECRET", slackSigningSecret)
			}
			if tt.args.envVarMisMatch {
				defer os.Unsetenv("SLACK_SIGNING_SECRET")
				os.Setenv("SLACK_SIGNING_SECRET", "somethingSilly")
			}
			if err := VerifySignature(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("VerifySignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
