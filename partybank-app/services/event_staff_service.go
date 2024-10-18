package services

import (
	"errors"
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
)

type EventStaffService interface {
	Create(createUserRequest *request.CreateEventStaffRequest) (*response.RaveResponse[string], error)
}

type raveEventStaffService struct {
	repositories.EventStaffRepository
	repositories.EventRepository
}

func NewEventStaffService(eventStaffRepository repositories.EventStaffRepository, eventRepository repositories.EventRepository) EventStaffService {
	return &raveEventStaffService{
		eventStaffRepository,
		eventRepository,
	}
}

func (eventStaffService *raveEventStaffService) Create(createUserRequest *request.CreateEventStaffRequest) (*response.RaveResponse[string], error) {
	event, err := eventStaffService.EventRepository.FindById(createUserRequest.EventId)
	if err != nil {
		return nil, errors.New("event not found")
	}
	mailService := NewMailService()
	for _, email := range createUserRequest.StaffEmails {
		savedStaff, err := updateEvent(email, event, eventStaffService.EventStaffRepository)
		if err != nil {
			return nil, err
		} else if savedStaff != nil {
			notification := request.NewEmailNotificationRequest(email, `
				welcome to rave, sign in using this email address
			`)
			_, err := mailService.Send(notification)
			if err != nil {
				return nil, err
			}
		}
	}

	return &response.RaveResponse[string]{Data: "event staffs invited"}, nil
}

func updateEvent(email string, event *models.Event, repo repositories.EventStaffRepository) (*models.EventStaff, error) {
	eventStaff := mapMailToEventStaff(email)
	eventStaff.EventID = event.ID
	event.EventStaff = append(event.EventStaff, eventStaff)
	savedStaff, err := repo.Save(eventStaff)
	return savedStaff, err
}

func mapMailToEventStaff(email string) *models.EventStaff {
	return &models.EventStaff{
		User: &models.User{
			Username: email,
			Role:     models.EVENT_STAFF,
		},
	}
}
