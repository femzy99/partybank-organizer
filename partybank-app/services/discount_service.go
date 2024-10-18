package services

import (
	"errors"
	"fmt"
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"gopkg.in/jeevatkm/go-model.v1"
)

type DiscountService interface {
	CreateDiscount(request *request.CreateDiscountRequest) (*response.CreateDiscountResponse, error)
}

type raveDiscountService struct {
	repositories.DiscountRepository
	TicketService
}

func NewDiscountService(discountRepository repositories.DiscountRepository, ticketService TicketService) DiscountService {
	return &raveDiscountService{
		discountRepository,
		ticketService,
	}
}

func (raveDiscountService *raveDiscountService) CreateDiscount(request *request.CreateDiscountRequest) (*response.CreateDiscountResponse, error) {
	ticketId := request.TicketId
	ticket, err := raveDiscountService.TicketService.GetTicketById(ticketId)
	isTicketNotExists := err != nil || ticket == nil
	if isTicketNotExists {
		return nil, errors.New(fmt.Sprintf("Cannot create Discount\n reason: Failed to find ticket with id: %d", ticketId))
	}
	discount := &models.Discount{}
	createDiscountResponse := &response.CreateDiscountResponse{}
	model.Copy(discount, request)
	discount.Ticket = ticket
	savedDiscount, err := raveDiscountService.Save(discount)
	if err != nil {
		return nil, err
	}
	model.Copy(createDiscountResponse, savedDiscount)
	return createDiscountResponse, nil
}
