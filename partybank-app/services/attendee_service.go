package services

import (
	"bytes"
	"errors"
	"github.com/djfemz/organizer-service/partybank-app/security"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"html/template"

	dtos "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/mappers"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"log"
)

type AttendeeService interface {
	Register(createAttendeeRequest *dtos.CreateAttendeeRequest) (*response.AttendeeResponse, error)
	GetAttendeeByUsername(username string) (*response.AttendeeResponse, error)
	FindAttendeeByUsername(username string) (*models.Attendee, error)
	UpdateAttendee(email string, updateAttendeeRequest *dtos.UpdateAttendeeRequest) (*response.AttendeeResponse, error)
}

type raveAttendeeService struct {
	repositories.AttendeeRepository
	MailService
}

func NewAttendeeService(attendeeRepository repositories.AttendeeRepository, mailService MailService) AttendeeService {
	return &raveAttendeeService{
		attendeeRepository,
		mailService,
	}
}

func (raveAttendeeService *raveAttendeeService) Register(createAttendeeRequest *dtos.CreateAttendeeRequest) (*response.AttendeeResponse, error) {
	attendee := mappers.MapCreateAttendeeRequestToAttendee(createAttendeeRequest)
	attendee, err := raveAttendeeService.Save(attendee)
	if err != nil {
		log.Println("Error: ", err.Error())
		return nil, errors.New("failed to create attendee service")
	}
	attendeeWelcomeMailRequest, err := buildNewAttendeeMessageFor(attendee)
	if err != nil {
		return nil, err
	}
	go raveAttendeeService.MailService.Send(attendeeWelcomeMailRequest)

	return mappers.MapAttendeeToAttendeeResponse(attendee), nil

}

func (raveAttendeeService *raveAttendeeService) GetAttendeeByUsername(username string) (*response.AttendeeResponse, error) {
	attendee, err := raveAttendeeService.FindByUsername(username)
	if err != nil {
		return nil, errors.New("failed to find attendee")
	}
	return mappers.MapAttendeeToAttendeeResponse(attendee), nil
}

func (raveAttendeeService *raveAttendeeService) FindAttendeeByUsername(username string) (*models.Attendee, error) {
	attendee, err := raveAttendeeService.FindByUsername(username)
	if err != nil {
		return nil, errors.New("attendee not found")
	}
	return attendee, nil
}

func (raveAttendeeService *raveAttendeeService) UpdateAttendee(username string, updateAttendeeRequest *dtos.UpdateAttendeeRequest) (*response.AttendeeResponse, error) {
	attendee, err := raveAttendeeService.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user with id not found")
	}
	attendee.FirstName = updateAttendeeRequest.FirstName
	attendee.LastName = updateAttendeeRequest.LastName
	attendee.PhoneNumber = updateAttendeeRequest.PhoneNumber
	attendee, err = raveAttendeeService.Save(attendee)
	if err != nil {
		return nil, errors.New("could not update user account")
	}
	res := &response.AttendeeResponse{
		Username: attendee.Username,
		Message:  "account updated successfully",
	}
	return res, nil
}

func buildNewAttendeeMessageFor(attendee *models.Attendee) (*dtos.EmailNotificationRequest, error) {
	tmpl, err := getAttendeeEmailTemplate(attendee)
	if err != nil {
		return nil, errors.New("could not get mail template")
	}
	return &dtos.EmailNotificationRequest{
		Sender: dtos.Sender{
			Name:  utils.APP_NAME,
			Email: utils.APP_EMAIL,
		},
		Recipients: []dtos.Recipient{
			{
				Name:  attendee.FirstName,
				Email: attendee.Username,
			},
		},
		Subject: "Welcome mail",
		Content: tmpl,
	}, nil
}

type attendeeMessage struct {
	Link string
}

func getAttendeeEmailTemplate(attendee *models.Attendee) (string, error) {
	token, err := security.GenerateAccessTokenFor(attendee)
	if err != nil {
		return "", err
	}
	message := &attendeeMessage{
		Link: "https://thepartybank.com/validate?" + "token=" + token,
	}
	mailTemplate, err := template.ParseFiles("rave-mail-template-new.html")
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}
	var body bytes.Buffer
	err = mailTemplate.Execute(&body, message)
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}
	return body.String(), nil
}
