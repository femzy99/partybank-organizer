package test

//
//import (
//	request "github.com/djfemz/rave/partybank-app/dtos/request"
//	"github.com/djfemz/rave/partybank-app/services"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//var addTicketRequest = &request.CreateTicketRequest{
//	Type:            "single",
//	Name:            "early birds",
//	NumberAvailable: 50,
//	Price:           5000.00,
//	EventID:         2,
//}
//
//var ticketService = services.NewTicketService()
//
//func TestAddTicket(t *testing.T) {
//	res, err := ticketService.CreateTicketFor(addTicketRequest)
//	assert.Nil(t, err)
//	assert.NotNil(t, res)
//}
//
//func TestGetTicketById(t *testing.T) {
//	ticket, err := ticketService.GetById(2)
//	assert.Nil(t, err)
//	assert.NotNil(t, ticket)
//}
//
//func TestGetAllTicketsForEvent(t *testing.T) {
//	tickets, err := ticketService.GetAllTicketsFor(2)
//	assert.Nil(t, err)
//	assert.NotNil(t, tickets)
//	assert.NotEmpty(t, tickets)
//}
