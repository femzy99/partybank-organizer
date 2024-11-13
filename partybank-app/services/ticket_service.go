package services

import (
	"errors"
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/mappers"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"gopkg.in/jeevatkm/go-model.v1"
	"log"
)

type TicketService interface {
	CreateTicketFor(eventId uint64, request *request.CreateTicketRequest) (addTicketResponse *models.Ticket, err error)
	AddTickets(eventId uint64, tickets []*request.CreateTicketRequest) ([]*response.TicketResponse, error)
	GetById(id uint64) (*response.TicketResponse, error)
	GetTicketById(id uint64) (*models.Ticket, error)
	GetAllTicketsFor(eventId uint64, pageNumber, pageSize int) ([]*response.TicketResponse, error)
	UpdateTicketSoldOutBy(reference string) (*response.TicketResponse, error)
	UpdateTicket(ticketId uint64, updateTicketRequest *request.UpdateTicketRequest) (*response.TicketResponse, error)
	EditTicket(ticketId uint64, editTicketRequest *request.EditTicketRequest) (editTicketResponse *response.TicketResponse, err error)
	EditTickets(eventId uint64, editTicketRequests []*request.EditTicketRequest) (editTicketResponses []*response.TicketResponse, err error)
}

type raveTicketService struct {
	repositories.TicketRepository
	EventService
}

func NewTicketService(ticketRepository repositories.TicketRepository, eventService EventService) TicketService {
	return &raveTicketService{
		ticketRepository,
		eventService,
	}
}

func (raveTicketService *raveTicketService) CreateTicketFor(eventId uint64, request *request.CreateTicketRequest) (ticket *models.Ticket, err error) {
	if request.ID > 0 {
		return nil, nil
	}
	event, err := raveTicketService.GetEventBy(eventId)
	if err != nil {
		log.Println("event: ", event)
		return nil, errors.New("event not found")
	}
	if utils.ExistsWithTicketName(event, request.Name) {
		return nil, errors.New("ticket already exists")
	}
	ticket = &models.Ticket{}
	errs := model.Copy(ticket, request)
	if len(errs) != 0 {
		log.Println(errs)
		return nil, errors.New("failed to create ticket")
	}
	ticket.Reference = utils.GenerateTicketReference()
	ticket.EventID = event.ID
	ticket.TicketPerks = request.TicketPerks
	ticket.ActivePeriod = &models.ActivePeriod{
		StartDate: request.SalesStartDate,
		EndDate:   request.SaleEndDate,
		StartTime: request.SalesStartTime,
		EndTime:   request.SalesEndTime,
	}
	ticket.Category = request.Category
	ticket.GroupTicketCapacity = request.GroupTicketCapacity
	savedTicket, err := raveTicketService.TicketRepository.Save(ticket)
	if err != nil {
		log.Println("error: ticket saving failed", err)
		return nil, errors.New("failed to save ticket")
	}
	event.Tickets = append(event.Tickets, savedTicket)
	err = raveTicketService.UpdateEvent(event)
	if err != nil {
		return nil, errors.New("failed to save ticket")
	}

	return savedTicket, nil
}

func (raveTicketService *raveTicketService) GetById(id uint64) (*response.TicketResponse, error) {
	ticket, err := raveTicketService.FindById(id)
	if err != nil {
		return nil, err
	}
	res := &response.TicketResponse{}
	errs := model.Copy(res, ticket)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	res.TicketPerks = ticket.TicketPerks
	return res, nil
}

func (raveTicketService *raveTicketService) GetTicketById(id uint64) (*models.Ticket, error) {
	ticket, err := raveTicketService.FindById(id)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (raveTicketService *raveTicketService) UpdateTicketSoldOutBy(reference string) (*response.TicketResponse, error) {
	ticket, err := raveTicketService.TicketRepository.FindTicketByReference(reference)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	ticket.IsSoldOutTicket = true
	ticket, err = raveTicketService.TicketRepository.Save(ticket)
	if err != nil {
		log.Println("Error: saving failed: ", err)
		return nil, errors.New(err.Error())
	}
	return mappers.MapTicketToTicketResponse(ticket), nil
}

func (raveTicketService *raveTicketService) GetAllTicketsFor(eventId uint64, pageNumber, pageSize int) ([]*response.TicketResponse, error) {
	tickets, err := raveTicketService.FindAllByEventId(eventId, pageNumber, pageSize)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return mappers.MapTicketsToTicketsResponse(tickets), nil
}

func (raveTicketService *raveTicketService) AddTickets(eventId uint64, tickets []*request.CreateTicketRequest) ([]*response.TicketResponse, error) {
	res := make([]*response.TicketResponse, 0)
	err := raveTicketService.DeleteTicketsFor(eventId)
	if err != nil {
		log.Println("error: ticket_service.go: 146 \t\tfailed to remove tickets from event")
	}
	if len(tickets) < 1 {
		return res, nil
	}
	for _, ticketRequest := range tickets {
		ticketResponse, _ := raveTicketService.CreateTicketFor(eventId, ticketRequest)
		ticket := mappers.MapTicketToTicketResponse(ticketResponse)
		res = append(res, ticket)
	}
	event, err := raveTicketService.EventService.GetEventBy(eventId)
	if err != nil {
		log.Println("error: ticket_service.go: 158 \t\tfailed to find event")
	}
	go utils.SendNewTicketMessageFor(event)
	return res, nil
}

func (raveTicketService *raveTicketService) UpdateTicket(ticketId uint64, updateTicketRequest *request.UpdateTicketRequest) (*response.TicketResponse, error) {
	foundTicket, err := raveTicketService.GetTicketById(ticketId)
	if err != nil {
		return nil, errors.New("failed to find ticket")
	}
	foundTicket.Name = updateTicketRequest.Name
	foundTicket.Type = updateTicketRequest.Type
	foundTicket.PurchaseLimit = updateTicketRequest.PurchaseLimit
	foundTicket.Capacity = updateTicketRequest.Capacity
	foundTicket.Colour = updateTicketRequest.Colour
	foundTicket.Price = updateTicketRequest.Price
	foundTicket.Stock = updateTicketRequest.Stock
	foundTicket.IsTransferPaymentFeesToGuest = updateTicketRequest.IsTransferPaymentFeesToGuest
	foundTicket.Category = updateTicketRequest.Category
	foundTicket.GroupTicketCapacity = updateTicketRequest.GroupTicketCapacity
	if foundTicket.ActivePeriod != nil {
		foundTicket.ActivePeriod.StartTime = updateTicketRequest.SalesStartTime
		foundTicket.ActivePeriod.EndTime = updateTicketRequest.SalesEndTime
		foundTicket.ActivePeriod.StartDate = updateTicketRequest.SalesStartDate
		foundTicket.ActivePeriod.EndDate = updateTicketRequest.SaleEndDate
	}
	foundTicket, err = raveTicketService.Save(foundTicket)
	if err != nil {
		return nil, errors.New("failed to save ticket")
	}
	event, err := raveTicketService.EventService.GetEventBy(foundTicket.EventID)
	go utils.SendNewTicketMessageFor(event)
	return mappers.MapTicketToTicketResponse(foundTicket), nil
}

func (raveTicketService *raveTicketService) EditTicket(ticketId uint64, editTicketRequest *request.EditTicketRequest) (editTicketResponse *response.TicketResponse, err error) {
	ticket, err := raveTicketService.GetTicketById(ticketId)
	if err != nil {
		return nil, errors.New("failed to find ticket")
	}
	ticket = mappers.MapEditTicketRequestToTicket(editTicketRequest, ticket)
	ticket, err = raveTicketService.TicketRepository.Save(ticket)
	if err != nil {
		log.Println("Error saving ticket in edit: ", err)
		return nil, errors.New("failed to edit ticket")
	}
	ticketResponse := mappers.MapTicketToTicketResponse(ticket)
	return ticketResponse, nil
}

func (raveTicketService *raveTicketService) EditTickets(eventId uint64, editTicketRequests []*request.EditTicketRequest) (editTicketResponses []*response.TicketResponse, err error) {
	tickets := make([]*models.Ticket, 0)
	if len(editTicketRequests) < 1 {
		err := raveTicketService.TicketRepository.DeleteTicketsFor(eventId)
		return nil, err
	}
	for _, ticketRequest := range editTicketRequests {
		if ticketRequest.ID == 0 {
			createTicketRequest := mappers.MapEditTicketToCreateTicket(ticketRequest)
			ticket, err := raveTicketService.CreateTicketFor(eventId, createTicketRequest)
			if err != nil {
				log.Println("error in edit tickets: ", err)
				return nil, errors.New("failed to add new ticket")
			}
			tickets = append(tickets, ticket)
			continue
		}
		ticket, err := raveTicketService.GetTicketById(ticketRequest.ID)
		if err != nil {
			log.Println("ERROR: ", err)
			return nil, errors.New("failed to find ticket")
		}
		ticket = mappers.MapEditTicketRequestToTicket(ticketRequest, ticket)
		tickets = append(tickets, ticket)
	}
	err = raveTicketService.DeleteAllNotIn(eventId, tickets)
	if err != nil {
		log.Println("Error removing tickets in edit: ", err)
		return nil, errors.New("failed to edit ticket")
	}
	err = raveTicketService.TicketRepository.SaveAll(tickets)
	if err != nil {
		log.Println("Error saving ticket in edit: ", err)
		return nil, errors.New("failed to edit ticket")
	}
	ticks := mappers.MapTicketsToTicketsResponse(tickets)
	return ticks, nil
}
