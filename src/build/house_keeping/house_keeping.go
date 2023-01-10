package house_keeping

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/riweston/acloudguru-client-go"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

var apiKey string
var consumerId string
var daysCap int
var licenseCap int

// init functions fail fast if the environment variables are not set

func init() {
	apiKey = os.Getenv("ACLOUDGURU_API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY not set")
		os.Exit(1)
	}
	consumerId = os.Getenv("ACLOUDGURU_CONSUMER_ID")
	if consumerId == "" {
		fmt.Println("CONSUMER_ID not set")
		os.Exit(1)
	}
}

func init() {
	var err error
	daysCap, err = GetDaysCap()
	if err != nil {
		fmt.Println("DAYS_CAP not set")
		os.Exit(1)
	}
	licenseCap, err = GetLicenseCap()
	if err != nil {
		fmt.Println("LICENSE_CAP not set")
		os.Exit(1)
	}
}

func init() {
	functions.CloudEvent("HouseKeeping", houseKeeping)
}

func houseKeeping(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	client, err := acg.NewClient(&apiKey, &consumerId)
	if err != nil {
		fmt.Errorf("error creating client: %v", err)
	}
	usersAll, err := client.GetUsersAll()
	if err != nil {
		fmt.Errorf("error getting users: %v", err)
	}
	usersProcessed := NewUsers(usersAll)
	usersDeactivate := usersProcessed.GetUsersToDeactivate()

	if len(usersDeactivate) > 0 {
		for _, user := range usersDeactivate {
			fmt.Println(user.Status, user.LastSeenTimestamp, user.Name)
			request := Request{
				emailAddress: user.Email,
				responseUrl:  msg.Message.Attributes["response_url"],
				requestType:  "deactivate",
			}
			request.PublishMessage()
		}
	}
	return nil
}

func (r *Request) PublishMessage() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic(os.Getenv("TOPIC_NAME"))
	defer topic.Stop()
	topic.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"user_email":   r.emailAddress,
			"response_url": r.responseUrl,
			"request_type": r.requestType,
		},
	})
}

func NewUsers(AllUsers *[]acg.User) *Users {
	users := new(Users)
	users.AllUsers = AllUsers
	users.LicenseCap = licenseCap
	users.DaysCap = daysCap
	users.FilterActiveUsers()
	users.FilterOldUsers()
	return users
}

// models

type Request struct {
	requestType  string
	responseUrl  string
	emailAddress string
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

type Users struct {
	LicenseCap  int
	DaysCap     int
	AllUsers    *[]acg.User
	ActiveUsers []acg.User
	OldUsers    []acg.User
}

func (r *Users) RemoveUsersFromSlice(users *[]acg.User, usersToRemove *[]acg.User) []acg.User {
	for _, user := range *usersToRemove {
		for i, user2 := range *users {
			if user.UserId == user2.UserId {
				*users = append((*users)[:i], (*users)[i+1:]...)
				break
			}
		}
	}
	return *users
}

func (r *Users) GetUsersToDeactivate() (result []acg.User) {
	ActiveUsers := r.ActiveUsers
	// If there are users older than the days cap add them to be processed and remove them from our working slice
	if len(r.OldUsers) > 0 {
		result = append(result, r.OldUsers...)
		ActiveUsers = r.RemoveUsersFromSlice(&r.ActiveUsers, &r.OldUsers)
	}
	// If there are still more users than the license cap add the remaining oldest users to be processed
	CheckCap := len(ActiveUsers) - len(result)
	if CheckCap > r.LicenseCap {
		for i := 0; i <= CheckCap-r.LicenseCap; i++ {
			element := (len(ActiveUsers) - 1) - i
			result = append(result, (ActiveUsers)[element])
		}
	}
	return
}

func GetLicenseCap() (int, error) {
	daysCapStr := os.Getenv("LICENSE_CAP")
	envVar, err := strconv.Atoi(daysCapStr)
	if err != nil {
		return envVar, err
	}
	return envVar, nil
}

func GetDaysCap() (int, error) {
	daysCapStr := os.Getenv("DAYS_CAP")
	envVar, err := strconv.Atoi(daysCapStr)
	if err != nil {
		return envVar, err
	}
	return envVar, nil
}

func (r *Users) FilterActiveUsers() {
	r.ActiveUsers = filterSlice(r.AllUsers, filterActiveUsers)
	sort.Slice(r.ActiveUsers, func(i, j int) bool {
		return (r.ActiveUsers)[i].LastSeenTimestamp.After((r.ActiveUsers)[j].LastSeenTimestamp)
	})
}

func (r *Users) FilterOldUsers() {
	r.OldUsers = filterSlice(&r.ActiveUsers, r.filterOldUsers)
}

// Helper function to filter a slice

func filterSlice(s *[]acg.User, t func(acg.User) bool) (ret []acg.User) {
	for _, v := range *s {
		if t(v) {
			ret = append(ret, v)
		}
	}
	return
}

func filterActiveUsers(user acg.User) bool {
	return user.Status == "Active"
}

func (r *Users) filterOldUsers(user acg.User) bool {
	return user.LastSeenTimestamp.Before(
		time.Now().AddDate(0, 0, -r.DaysCap))
}
