package slack

import (
	"net/http"
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
