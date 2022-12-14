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

func VerifyTimeStamp(r *http.Request) (string, error) {
	headerTimeStamp := r.Header["x-slack-request-timestamp"][0]
	if len(headerTimeStamp) == 0 {
		log.Println("x-slack-request-timestamp header is missing")
		return "", fmt.Errorf("x-slack-request-timestamp header is missing")
	}
	headerTimeStampInt, err := strconv.ParseInt(headerTimeStamp, 10, 64)
	if err != nil {
		log.Println("x-slack-request-timestamp header value is invalid")
		return "", fmt.Errorf("x-slack-request-timestamp header value is invalid")
	}
	currentTime := time.Now().Unix()
	if (currentTime - headerTimeStampInt) > (60 * 5) {
		log.Println("x-slack-request-timestamp is too old")
		return "", fmt.Errorf("x-slack-request-timestamp is too old")
	}
	return headerTimeStamp, nil
}

func GenerateSigBaseString(r *http.Request) string {
	headerTimeStamp, err := VerifyTimeStamp(r)
	if err != nil {
		err.Error()
	}
	requestBody, _ := io.ReadAll(r.Body)
	return fmt.Sprint("v0:" + headerTimeStamp + ":" + string(requestBody))
}

func HashSigBaseString(sigBaseString string) (string, error) {
	slackSigningSecret := os.Getenv("SLACK_SIGNING_SECRET")
	if slackSigningSecret == "" {
		log.Println("SLACK_SIGNING_SECRET is missing")
		return "", fmt.Errorf("SLACK_SIGNING_SECRET is not set")
	}
	h := hmac.New(sha256.New, []byte(slackSigningSecret))
	_, err := h.Write([]byte(sigBaseString))
	if err != nil {
		log.Println("Failed to hash the signature base string")
		return "", fmt.Errorf("error hashing sigBaseString: %v", err)
	}
	return fmt.Sprintf("v0=%x", h.Sum(nil)), nil
}

func VerifySignature(r *http.Request) error {
	sigBaseString := GenerateSigBaseString(r)
	slackSignature := r.Header["x-slack-signature"][0]
	if len(slackSignature) == 0 {
		log.Println("x-slack-signature header is missing")
		return fmt.Errorf("x-slack-signature header is missing")
	}
	hashedSigBaseString, err := HashSigBaseString(sigBaseString)
	if err != nil {
		log.Println("Unknown Error: Failed to hash the signature base string")
		return err
	}
	if slackSignature != hashedSigBaseString {
		log.Println("x-slack-signature is mismatch")
		return fmt.Errorf("x-slack-signature is mismatch")
	}
	return nil
}
