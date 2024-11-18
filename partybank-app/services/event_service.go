package services

import (
	"bytes"
	"errors"
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/mappers"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type EventService interface {
	Create(createEventRequest *request.CreateEventRequest) (*response.EventResponse, error)
	GetById(id uint64) (*response.EventResponse, error)
	GetEventBy(id uint64) (*models.Event, error)
	GetEventByReference(reference string) (*response.EventResponse, error)
	DiscoverEvents(page int, size int) ([]*response.EventResponse, error)
	UpdateEventInformation(id uint64, updateRequest *request.UpdateEventRequest) (*response.EventResponse, error)
	UpdateEvent(event *models.Event) error
	GetAllEventsFor(organizerId uint64, pageNumber int, pageSize int) ([]*response.EventResponse, error)
	PublishEvent(eventId uint64) (*response.EventResponse, error)
	GetOrganizerFor(eventId uint64) (uint64, error)
	GetAllEventsForOrganizer(organizerId uint64, page, size int) ([]*response.EventResponse, error)
	SetTicketService(service TicketService)
	DeleteEventBy(eventId uint64) (string, error)
}

type raveEventService struct {
	repositories.EventRepository
	OrganizerService
	SeriesService
	TicketService
}

func NewEventService(eventRepository repositories.EventRepository,
	organizerService OrganizerService,
	seriesService SeriesService,
	ticketService TicketService,
) EventService {
	return &raveEventService{
		eventRepository,
		organizerService,
		seriesService,
		ticketService,
	}
}

func (raveEventService *raveEventService) Create(createEventRequest *request.CreateEventRequest) (*response.EventResponse, error) {
	event := mapCreateEventRequestToEvent(createEventRequest)
	var err error
	var calendar *models.Series
	org, err := raveEventService.OrganizerService.GetById(createEventRequest.OrganizerId)
	if err != nil || org == nil {
		log.Println("err finding organizer: ", err)
		return nil, err
	}
	if createEventRequest.SeriesId == 0 {
		calendar, err = raveEventService.GetPublicCalendarFor(createEventRequest.OrganizerId)
		if err != nil {
			log.Println("error finding public calendar: ", err)
			return nil, err
		}
		log.Println("found public calendar: ", calendar)
	} else {
		calendar, err = raveEventService.SeriesService.GetById(createEventRequest.SeriesId)
		if err != nil {
			return nil, err
		}
	}
	updateEventDetails(createEventRequest, event, calendar)
	savedEvent, err := raveEventService.Save(event)
	if err != nil {
		log.Println("err saving event: ", err)
		return nil, err
	}
	_, err = raveEventService.AddEventToSeries(calendar.ID, savedEvent)
	if err != nil {
		return nil, err
	}
	if len(createEventRequest.Tickets) > 0 {
		_, err = raveEventService.TicketService.AddTickets(savedEvent.ID, createEventRequest.Tickets)
		if err != nil {
			log.Println("error adding tickets to event")
		}
	}
	res := mappers.MapEventToEventResponse("event created successfully", savedEvent)
	res.SeriesName = calendar.Name
	return res, nil
}

func updateEventDetails(createEventRequest *request.CreateEventRequest, event *models.Event, calendar *models.Series) {
	event.SeriesID = calendar.ID
	event.CreatedBy = calendar.Name
	event.PublicationState = models.DRAFT
	event.IsNotificationEnabled = createEventRequest.IsNotificationEnabled
	event.CreatedBy = strconv.Itoa(int(createEventRequest.OrganizerId))
	event.Location = &models.Location{
		Longitude: createEventRequest.Longitude,
		Latitude:  createEventRequest.Latitude,
		City:      createEventRequest.City,
		State:     createEventRequest.State,
		Country:   createEventRequest.Country,
	}
}

func (raveEventService *raveEventService) GetById(id uint64) (*response.EventResponse, error) {
	foundEvent, err := raveEventService.FindById(id)
	if err != nil {
		return nil, err
	}
	log.Println("event: ", *foundEvent)
	series, err := raveEventService.SeriesService.GetById(foundEvent.SeriesID)
	if err != nil {
		return nil, errors.New("failed to find series")
	}
	eventResponse := mappers.MapEventToEventResponse("Success", foundEvent)
	eventResponse.SeriesName = series.Name
	return eventResponse, nil
}

func (raveEventService *raveEventService) GetEventBy(id uint64) (*models.Event, error) {
	return raveEventService.FindById(id)
}

func (raveEventService *raveEventService) UpdateEventInformation(id uint64, updateRequest *request.UpdateEventRequest) (*response.EventResponse, error) {
	foundEvent, err := raveEventService.FindById(id)
	if err != nil {
		return nil, err
	}
	foundEvent = mappers.MapUpdateEventRequestToEvent(updateRequest, foundEvent)
	savedEvent, err := raveEventService.Save(foundEvent)
	if err != nil {
		return nil, errors.New("failed to save event")
	}
	_, err = raveEventService.TicketService.EditTickets(id, updateRequest.Tickets)
	if err != nil {
		log.Println("Error adding ticket: ", err.Error())
		return nil, errors.New("failed to update event")
	}
	savedEvent, err = raveEventService.Save(savedEvent)
	if err != nil {
		return nil, errors.New("failed to update event during event update")
	}
	go utils.SendNewTicketMessageFor(savedEvent)
	updateEventResponse := mappers.MapEventToEventResponse("event update successful", savedEvent)
	return updateEventResponse, nil
}

func (raveEventService *raveEventService) UpdateEvent(event *models.Event) error {
	_, err := raveEventService.Save(event)
	return err
}

func (raveEventService *raveEventService) GetAllEventsFor(calendarId uint64, pageNumber, pageSize int) ([]*response.EventResponse, error) {
	events, err := raveEventService.FindAllByCalendar(calendarId, pageNumber, pageSize)
	eventsResponses := make([]*response.EventResponse, 0)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eventResponse := mappers.MapEventToEventResponse("success", event)
		eventsResponses = append(eventsResponses, eventResponse)
	}
	return eventsResponses, nil
}

func (raveEventService *raveEventService) DiscoverEvents(page int, size int) ([]*response.EventResponse, error) {
	events, err := raveEventService.FindAllPublishedByPage(page, size)
	if err != nil {
		log.Println("error: ", err.Error())
		return make([]*response.EventResponse, 0), err
	}
	log.Println("events: ", events)
	if events != nil && len(events) > 0 {
		series, err := raveEventService.SeriesService.GetById(events[0].SeriesID)
		if err != nil {
			return make([]*response.EventResponse, 0), err
		}
		allEvents := mappers.MapEventsToEventResponses(events, series)
		return allEvents, nil
	}

	return make([]*response.EventResponse, 0), errors.New("no events found")
}

func (raveEventService *raveEventService) GetEventByReference(reference string) (*response.EventResponse, error) {
	event, err := raveEventService.EventRepository.FindByReference(reference)
	if err != nil {
		return nil, errors.New("failed to find requested event")
	}
	series, err := raveEventService.SeriesService.GetById(event.SeriesID)
	if err != nil {
		return nil, errors.New("failed to find series")
	}
	eventResponse := mappers.MapEventToEventResponse("Success", event)
	eventResponse.SeriesName = series.Name
	return eventResponse, nil
}

func (raveEventService *raveEventService) PublishEvent(eventId uint64) (*response.EventResponse, error) {
	event, err := raveEventService.GetEventBy(eventId)
	if err != nil {
		return nil, errors.New("event not found")
	}
	if event.Tickets != nil && len(event.Tickets) > 0 {
		event.PublicationState = models.PUBLISHED
	} else if event.Tickets == nil || len(event.Tickets) <= 0 {
		return nil, errors.New("failed to publish event without tickets, add ticket to event first")
	}
	event, err = raveEventService.Save(event)
	if err != nil {
		return nil, errors.New("failed to save event")
	}
	series, err := raveEventService.SeriesService.GetById(event.SeriesID)
	if err != nil {
		return nil, errors.New("failed to find series")
	}
	eventResponse := mappers.MapEventToEventResponse("Success", event)
	eventResponse.SeriesName = series.Name
	return eventResponse, nil
}

func (raveEventService *raveEventService) GetOrganizerFor(SeriesId uint64) (uint64, error) {
	organizerId, err := raveEventService.SeriesService.GetSeriesOrganizer(SeriesId)
	if err != nil {
		return 0, errors.New("failed to find event organizer")
	}
	return organizerId, nil
}

func (raveEventService *raveEventService) SetTicketService(service TicketService) {
	raveEventService.TicketService = service
}

func (raveEventService *raveEventService) GetAllEventsForOrganizer(organizerId uint64, page, size int) ([]*response.EventResponse, error) {
	events, err := raveEventService.FindAllByOrganizer(organizerId, page, size)
	if err != nil {
		return nil, errors.New("organizer not found")
	}
	eventsResponses := make([]*response.EventResponse, 0)
	for _, event := range events {
		eventRes := mappers.MapEventToEventResponse("Success", event)
		series, err := raveEventService.SeriesService.GetById(event.SeriesID)
		if err != nil {
			return nil, errors.New("series not found")
		}
		eventRes.SeriesName = series.Name
		eventsResponses = append(eventsResponses, eventRes)
	}
	sort.Slice(eventsResponses, func(currentIndex, nextIndex int) bool {
		return eventsResponses[currentIndex].ID > eventsResponses[nextIndex].ID
	})
	return eventsResponses, err
}

func (raveEventService *raveEventService) DeleteEventBy(eventId uint64) (string, error) {
	event, err := raveEventService.FindById(eventId)
	if err != nil {
		return "", errors.New("failed to find event")
	}
	event.PublicationState = "ARCHIVED"
	err = raveEventService.EventRepository.DeleteById(eventId)
	if err != nil {
		return "", errors.New("failed to delete event")
	}
	paymentServiceDeleteEndpoint := os.Getenv("DELETE_EVENT_ENDPOINT_PAYMENT_SERVICE")
	paymentServiceDeleteEndpoint = paymentServiceDeleteEndpoint + event.Reference
	req, err := http.NewRequest(http.MethodPost, paymentServiceDeleteEndpoint, bytes.NewReader([]byte("")))

	if err != nil {
		log.Println("ERROR: failed to create delete request to payment side")
		return "", errors.New("failed to delete event")
	}
	client := &http.Client{}
	res, err := client.Do(req)
	log.Println("delete event response from payment service: ", res.StatusCode)
	if err != nil || res.StatusCode != 200 {
		log.Println("ERROR: failed to send delete request to payment side", err)
		return "", errors.New("failed to delete event")
	}
	return "event deleted successfully", nil
}

func mapCreateEventRequestToEvent(createEventRequest *request.CreateEventRequest) *models.Event {
	return &models.Event{
		Name: createEventRequest.Name,
		Location: &models.Location{
			Longitude: createEventRequest.Longitude,
			Latitude:  createEventRequest.Latitude,
			Address:   createEventRequest.Address,
		},
		EventDate:             createEventRequest.Date,
		StartTime:             createEventRequest.StartTime,
		EndTime:               createEventRequest.EndTime,
		SeriesID:              createEventRequest.SeriesId,
		IsNotificationEnabled: createEventRequest.IsNotificationEnabled,
		ContactInformation:    createEventRequest.ContactInformation,
		Description:           createEventRequest.Description,
		Status:                models.UPCOMING,
		EventTheme:            createEventRequest.EventTheme,
		AttendeeTerm:          createEventRequest.AttendeeTerm,
		Venue:                 createEventRequest.Venue,
		ImageUrl:              createEventRequest.ImageUrl,
		Reference:             utils.GenerateEventReference(),
	}
}
