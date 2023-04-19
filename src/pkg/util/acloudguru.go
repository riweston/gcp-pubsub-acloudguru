package util

import (
	acg "github.com/riweston/acloudguru-client-go"
	"sort"
	"time"
)

type Users struct {
	LicenseCap  int
	DaysCap     int
	AllUsers    *[]acg.User
	ActiveUsers []acg.User
	OldUsers    []acg.User
}

func NewUsers(allUsers *[]acg.User, licenseCap int, daysCap int) *Users {
	users := new(Users)
	users.AllUsers = allUsers
	users.LicenseCap = licenseCap
	users.DaysCap = daysCap
	users.FilterActiveUsers()
	users.FilterOldUsers()
	return users
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
