package test

import (
	dtos "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/services"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestUpdateAttendee(t *testing.T) {
	username := utils.GenerateTicketReference()
	attendee := &models.Attendee{
		FirstName: "John Doe",
		User: &models.User{
			Username: username,
			Role:     models.ATTENDEE,
		},
		PhoneNumber: "09087654356",
	}
	db := repositories.Connect()
	attendeeRepository := repositories.NewAttendeeRepository(db)
	attendeeService := services.NewAttendeeService(attendeeRepository, nil)
	attendee, _ = attendeeRepository.Save(attendee)
	updateAttendeeRequest := &dtos.UpdateAttendeeRequest{
		FullName:    "James Doe",
		PhoneNumber: "09087654356",
	}

	res, err := attendeeService.UpdateAttendee(username, updateAttendeeRequest)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	attendee, err = attendeeRepository.FindById(attendee.ID)
	log.Println("attendee: ", *attendee)
	assert.Nil(t, err)
	assert.NotNil(t, attendee)
	assert.Equal(t, username, attendee.Username)
	assert.Equal(t, "James Doe", attendee.FirstName)
}
