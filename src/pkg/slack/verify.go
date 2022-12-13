package slack

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type HttpRequest struct {
	UserId      string
	ResponseUrl string
}

func (h *HttpRequest) Verify() bool {
	return false
}

func VerifyTimeStamp(r *http.Request) bool {
	headerTimeStamp := r.Header.Get("X-Slack-Request-Timestamp")
	if headerTimeStamp == "" {
		log.Println("X-Slack-Request-Timestamp header is missing")
		return false
	}
	headerTimeStampInt, err := strconv.ParseInt(headerTimeStamp, 10, 64)
	if err != nil {
		log.Println("X-Slack-Request-Timestamp header value is invalid")
		return false
	}
	currentTime := time.Now().Unix()
	if (currentTime - headerTimeStampInt) > (60 * 5) {
		log.Println("X-Slack-Request-Timestamp is too old")
		return false
	}
	return true
}
