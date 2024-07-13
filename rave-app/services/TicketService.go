package services

import (
	"bytes"
	"encoding/json"
	"errors"
	request "github.com/djfemz/rave/rave-app/dtos/request"
	response "github.com/djfemz/rave/rave-app/dtos/response"
	"github.com/djfemz/rave/rave-app/models"
	"github.com/djfemz/rave/rave-app/repositories"
	"gopkg.in/jeevatkm/go-model.v1"
	"log"
	"net/http"
	"os"
)

type TicketService interface {
	CreateTicketFor(request *request.CreateTicketRequest) (addTicketResponse *response.TicketResponse, err error)
	GetTicketById(id uint64) (*response.TicketResponse, error)
	GetAllTicketsFor(eventId uint64) ([]*models.Ticket, error)
}

var ticketRepository = repositories.NewTicketRepository()

type raveTicketService struct {
}

func NewTicketService() TicketService {
	return &raveTicketService{}
}

func (raveTicketService *raveTicketService) CreateTicketFor(request *request.CreateTicketRequest) (addTicketResponse *response.TicketResponse, err error) {
	eventService := NewEventService()
	event, err := eventService.GetEventBy(request.EventId)

	ticket := &models.Ticket{}
	errs := model.Copy(ticket, request)
	if len(errs) != 0 {
		log.Println(errs)
		return nil, errors.New("failed to create ticket")
	}
	savedTicket, err := ticketRepository.Save(ticket)
	if err != nil {
		return nil, errors.New("failed to save ticket")
	}
	event.Tickets = append(event.Tickets, ticket)
	err = eventService.UpdateEvent(event)
	if err != nil {
		return nil, errors.New("failed to save ticket")
	}
	createTicketResponse := &response.TicketResponse{}
	errs = model.Copy(createTicketResponse, savedTicket)

	log.Println("new ticket created: ", savedTicket)
	go sendNewTicketEvent(event, createTicketResponse)
	return createTicketResponse, nil
}

func (raveTicketService *raveTicketService) GetTicketById(id uint64) (*response.TicketResponse, error) {
	ticket, err := ticketRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	res := &response.TicketResponse{}
	errs := model.Copy(res, ticket)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return res, nil
}

func (raveTicketService *raveTicketService) GetAllTicketsFor(eventId uint64) ([]*models.Ticket, error) {
	tickets, err := ticketRepository.FindAllByEventId(eventId)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func sendNewTicketEvent(event *models.Event, ticketResponse *response.TicketResponse) {
	ticketMessage := buildTicketMessage(event, ticketResponse)
	body, err := json.Marshal(ticketMessage)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, os.Getenv("TICKET_SERVICE_URL"), bytes.NewReader(body))
	req.Header.Add("Content-Type", APPLICATION_JSON_VALUE)
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal("Error: ", err)
	}

}

func buildTicketMessage(event *models.Event, ticketResponse *response.TicketResponse) *request.NewTicketMessage {
	return &request.NewTicketMessage{
		Type:                       ticketResponse.Type,
		Name:                       ticketResponse.Name,
		Stock:                      ticketResponse.Stock,
		NumberAvailable:            ticketResponse.NumberAvailable,
		Price:                      ticketResponse.Price,
		DiscountCode:               ticketResponse.DiscountCode,
		DiscountPrice:              ticketResponse.DiscountPrice,
		PurchaseLimit:              ticketResponse.PurchaseLimit,
		Percentage:                 ticketResponse.Percentage,
		AvailableDiscountedTickets: ticketResponse.AvailableDiscountedTickets,
		EventName:                  event.Name,
		Description:                event.Description,
		Location:                   event.Location,
	}

}