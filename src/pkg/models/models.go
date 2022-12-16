package models

type Request struct {
	userId      string
	responseUrl string
}

type Activate struct {
	userEmail   string
	responseUrl string
	requestType string
}
