package test

//
//import (
//	request "github.com/djfemz/rave/partybank-app/dtos/request"
//	"github.com/djfemz/rave/partybank-app/services"
//	"github.com/stretchr/testify/assert"
//	"log"
//	"testing"
//)
//
//var eventService = services.NewEventService()
//
//func TestCreateEvent(t *testing.T) {
//	createEvent := &request.CreateEventRequest{
//		Name:               "test event",
//		Location:           "Sabo Yaba",
//		Date:               "2024-03-23",
//		StartTime:               "12:00:00",
//		ContactInformation: "09023456789",
//		Description:        "this is a test event",
//	}
//
//	res, err := eventService.Create(createEvent)
//	log.Println("event response: ", *res)
//	assert.NotNil(t, res)
//	assert.Nil(t, err)
//}
//
//func TestGetEventById(t *testing.T) {
//	event, err := eventService.GetById(1)
//	assert.NotNil(t, event)
//	assert.Nil(t, err)
//}
//
//func TestEditEventDetails(t *testing.T) {
//	updateRequest := &request.UpdateEventRequest{
//		Name:               "test event",
//		Location:           "Sabo Yaba",
//		Date:               "2024-03-23",
//		StartTime:               "12:00:00",
//		ContactInformation: "09023256887",
//		Description:        "this is a test event",
//	}
//
//	updateResponse, err := eventService.UpdateEventInformation(2, updateRequest)
//	assert.NotNil(t, updateResponse)
//	assert.Equal(t, updateResponse.ContactInformation, updateRequest.ContactInformation)
//	assert.Nil(t, err)
//}
//
//func TestGetAllEventsForOrganizer(t *testing.T) {
//	events, err := eventService.GetAllEventsFor(1)
//	assert.Nil(t, err)
//	assert.NotNil(t, events)
//	assert.NotEmpty(t, events)
//}
