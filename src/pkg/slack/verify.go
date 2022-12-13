package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func GenerateSigBaseString(r *http.Request) string {
	headerTimeStamp := r.Header.Get("X-Slack-Request-Timestamp")
	requestBody, _ := io.ReadAll(r.Body)
	return fmt.Sprint("v0:" + headerTimeStamp + ":" + string(requestBody))
}

func HashSigBaseString(sigBaseString string) (string, error) {
	slackSigningSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if slackSigningSecret == "" {
		return "", fmt.Errorf("SLACK_SIGNING_SECRET is not set")
	}
	h := hmac.New(sha256.New, []byte(slackSigningSecret))
	_, err := h.Write([]byte(sigBaseString))
	if err != nil {
		return "", fmt.Errorf("error hashing sigBaseString: %v", err)
	}
	return fmt.Sprintf("v0=%x", h.Sum(nil)), nil
}

func VerifySignature(r *http.Request) error {
	sigBaseString := GenerateSigBaseString(r)
	slackSignature := r.Header.Get("X-Slack-Signature")
	hashedSigBaseString, err := HashSigBaseString(sigBaseString)
	if err != nil {
		return err
	}
	if slackSignature != hashedSigBaseString {
		return fmt.Errorf("X-Slack-Signature is invalid")
	}
	return nil
}
