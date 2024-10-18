package services

import (
	"errors"
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/security/otp"

	"log"
)

type OrganizerService interface {
	Create(createOrganizerRequest *request.CreateUserRequest) (*response.CreateOrganizerResponse, error)
	GetByUsername(username string) (*models.Organizer, error)
	UpdateOtpFor(id uint64, testOtp *otp.OneTimePassword) (*models.Organizer, error)
	GetById(id uint64) (*models.Organizer, error)
	AddEventStaff(staff *request.AddEventStaffRequest) (*response.RaveResponse[string], error)
	GetByOtp(otp string) (*models.Organizer, error)
	IssueTicketTo(organizerId uint64, issueTicketRequest *request.IssueTicketRequest) (*response.RaveResponse[string], error)
}

type appOrganizerService struct {
	repository        repositories.OrganizerRepository
	eventStaffService EventStaffService
	seriesService     SeriesService
	ticketService     TicketService
	attendeeService   AttendeeService
}

func NewOrganizerService(organizerRepository repositories.OrganizerRepository,
	eventStaffService EventStaffService,
	seriesService SeriesService,
	ticketService TicketService,
	attendeeService AttendeeService) OrganizerService {
	return &appOrganizerService{
		repository:        organizerRepository,
		eventStaffService: eventStaffService,
		seriesService:     seriesService,
		ticketService:     ticketService,
		attendeeService:   attendeeService,
	}
}

func (organizerService *appOrganizerService) Create(createOrganizerRequest *request.CreateUserRequest) (*response.CreateOrganizerResponse, error) {
	organizer := mapCreateOrganizerRequestTo(createOrganizerRequest)
	password := otp.GenerateOtp()
	mailService := NewMailService()

	organizer.Otp = password
	savedOrganizer, err := organizerService.repository.Save(organizer)
	createCalendarRequest := &request.CreateSeriesRequest{
		Name:        "Public",
		OrganizerID: savedOrganizer.ID,
	}
	calendarResponse, err := organizerService.seriesService.AddSeries(createCalendarRequest)
	if err != nil {
		return nil, err
	}
	calendar, err := organizerService.seriesService.GetById(calendarResponse.ID)

	if err != nil {
		return nil, err
	}
	savedOrganizer.Series = append(savedOrganizer.Series, calendar)
	savedOrganizer, err = organizerService.repository.Save(savedOrganizer)

	go func() {
		mailService.Send(request.NewEmailNotificationRequest(CreateNewOrganizerEmail(password.Code), organizer.Username))
	}()
	if savedOrganizer != nil {
		return &response.CreateOrganizerResponse{
			Message:  response.USER_CREATED_SUCCESSFULLY,
			Username: savedOrganizer.Username,
		}, nil
	}
	return nil, err
}

func (organizerService *appOrganizerService) GetByUsername(username string) (*models.Organizer, error) {
	organizer, err := organizerService.repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return organizer, err
}

func (organizerService *appOrganizerService) UpdateOtpFor(id uint64, otp *otp.OneTimePassword) (*models.Organizer, error) {
	organizerRepository := organizerService.repository
	organizer, err := organizerRepository.FindById(id)
	if organizer != nil {
		organizer.Otp = otp
		organizer, err = organizerRepository.Save(organizer)
		if err != nil {
			return nil, err
		}
		return organizer, nil
	} else {
		return nil, err
	}
}

func (organizerService *appOrganizerService) GetById(id uint64) (*models.Organizer, error) {
	organizationRepository := organizerService.repository
	org, err := organizationRepository.FindById(id)
	if org == nil {
		return nil, err
	}
	return org, nil
}

func (organizerService *appOrganizerService) AddEventStaff(addStaffRequest *request.AddEventStaffRequest) (*response.RaveResponse[string], error) {
	res, err := organizerService.eventStaffService.Create(&request.CreateEventStaffRequest{StaffEmails: addStaffRequest.StaffEmails, EventId: addStaffRequest.EventId})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (organizerService *appOrganizerService) GetByOtp(otp string) (*models.Organizer, error) {
	organizerRepository := organizerService.repository
	log.Println("repo:", organizerRepository)
	return organizerRepository.FindByOtp(otp)
}

func (organizerService *appOrganizerService) IssueTicketTo(organizerId uint64, issueTicketRequest *request.IssueTicketRequest) (*response.RaveResponse[string], error) {
	organizer, err := organizerService.repository.FindById(organizerId)
	if err != nil {
		return nil, errors.New("organizer not found")
	}
	ticket, err := organizerService.ticketService.GetTicketById(issueTicketRequest.TicketId)
	if err != nil {
		return nil, errors.New("ticket not found")
	}
	attendee, err := organizerService.attendeeService.FindAttendeeByUsername(issueTicketRequest.AttendeeUsername)
	if err != nil {
		return nil, errors.New("attendee not found")
	}
	issuedTicket := &models.IssuedTicket{
		Issuer:   organizer,
		Attendee: attendee,
		Ticket:   ticket,
	}
	issuedTicketRepo := repositories.NewIssuedTicketRepository()
	issuedTicket, err = issuedTicketRepo.Save(issuedTicket)
	if err != nil {
		return nil, errors.New("failed to save issued ticket")
	}
	return &response.RaveResponse[string]{Data: "ticket has been issued to attendee"}, nil
}

func mapCreateOrganizerRequestTo(organizerRequest *request.CreateUserRequest) *models.Organizer {
	log.Println("organizerRequest", organizerRequest)
	return &models.Organizer{
		User: &models.User{
			Username: organizerRequest.Username,
			Role:     models.ORGANIZER,
		},
	}
}

func CreateNewOrganizerEmail(content string) string {
	return content
}
