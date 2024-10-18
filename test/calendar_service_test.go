package test

//
//
//import (
//	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
//	"github.com/djfemz/organizer-service/partybank-app/models"
//	"github.com/djfemz/organizer-service/partybank-app/services"
//	"github.com/stretchr/testify/assert"
//	"gopkg.in/jeevatkm/go-model.v1"
//	"log"
//	"testing"
//	"time"
//)
//
//var calendarService services.SeriesService
//
//func TestCreateCalendar(t *testing.T) {
//	calendarService = services.NewSeriesService()
//	createCalendarRequest := &request.CreateSeriesRequest{
//		Name:        "test",
//		OrganizerID: 1,
//		Description: "test desc",
//		ImageUrl:    "https://image.com",
//	}
//
//	response, err := calendarService.AddSeries(createCalendarRequest)
//	assert.NotNil(t, response)
//	assert.NotNil(t, response.Message)
//	assert.Nil(t, err)
//}
//
//func TestGetCalendar(t *testing.T) {
//	calendarService = services.NewSeriesService()
//	calendar, err := calendarService.GetById(1)
//	log.Println("calendar: ", calendar.Events)
//	assert.Nil(t, err)
//	assert.NotNil(t, calendar)
//}
//
//func TestAddEventToCalendar(t *testing.T) {
//	req := &request.CreateEventRequest{
//		SeriesId: 1,
//		Name:     "rave",
//		Address: "Abuja",
//
//		StartTime:     time.Now().String(),
//	}
//	eventService := services.NewEventService()
//	event, err := eventService.Create(req)
//	createdEvent := &models.Event{}
//	model.Copy(createdEvent, event)
//	assert.Nil(t, err)
//	calendarService = services.NewSeriesService()
//	log.Println("eve: ", createdEvent)
//	calendar, err := calendarService.AddEventToCalendar(1, createdEvent)
//
//	assert.Nil(t, err)
//	assert.NotNil(t, calendar)
//	assert.NotEmpty(t, calendar.Events)
//	assert.GreaterOrEqual(t, len(calendar.Events), 1)
//
//}
