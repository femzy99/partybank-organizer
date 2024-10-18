package services

import (
	"bytes"
	"encoding/json"
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"log"
	"net/http"
	"os"
)

const (
	API_KEY_VALUE          = "api-key"
	CONTENT_TYPE_KEY       = "Content-Type"
	APPLICATION_JSON_VALUE = "application/json"
	ACCEPT_HEADER_KEY      = "accept"
)

type MailService interface {
	Send(emailRequest *request.EmailNotificationRequest) (string, error)
}

type raveMailService struct{}

func NewMailService() MailService {
	return &raveMailService{}
}

func (raveMailService *raveMailService) Send(emailRequest *request.EmailNotificationRequest) (string, error) {
	jsonData, _ := json.Marshal(emailRequest)
	req, err := http.NewRequest(http.MethodPost, os.Getenv("MAIL_API_URL"), bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	addHeadersTo(req)

	client := &http.Client{}
	log.Println("Sending mail: ", req)
	if _, err = client.Do(req); err != nil {
		return "Mail Sending failed", err
	}
	return response.MAIL_SENDING_SUCCESS_MESSAGE, nil
}

func addHeadersTo(req *http.Request) {
	req.Header.Add(API_KEY_VALUE, os.Getenv("MAIL_API_KEY"))
	req.Header.Add(CONTENT_TYPE_KEY, APPLICATION_JSON_VALUE)
	req.Header.Add(ACCEPT_HEADER_KEY, APPLICATION_JSON_VALUE)
}
