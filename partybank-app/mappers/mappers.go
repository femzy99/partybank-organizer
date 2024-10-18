package mappers

import (
	"log"
	"time"

	dtos "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/models"
)

func MapSeriesCollectionToSeriesResponseCollection(series []*models.Series, organizer *models.Organizer) []*response.SeriesResponse {
	seriesResponses := make([]*response.SeriesResponse, 0)
	for _, userSeries := range series {
		seriesResponse := &response.SeriesResponse{
			ID:          userSeries.ID,
			OrganizerID: userSeries.OrganizerID,
			Name:        userSeries.Name,
			Description: userSeries.Description,
			ImageUrl:    userSeries.ImageUrl,
			Events:      MapEventsToEventResponses(userSeries.Events, series[0]),
		}
		seriesResponses = append(seriesResponses, seriesResponse)
	}
	return seriesResponses
}

func MapEventsToEventResponses(events []*models.Event, series *models.Series) []*response.EventResponse {
	responses := make([]*response.EventResponse, 0)
	for _, event := range events {
		ticketResponses := GetTicketsFrom(event)
		log.Println("tickets: ", ticketResponses)
		eventResponse := MapEventToEventResponse(event)

		eventResponse.Tickets = ticketResponses
		eventResponse.SeriesLogo = series.Logo
		responses = append(responses, eventResponse)
	}
	return responses
}

func GetTicketsFrom(event *models.Event) []*response.TicketResponse {
	log.Println("tickkets: ", event.Tickets)
	ticketResponses := make([]*response.TicketResponse, 0)
	for _, ticket := range event.Tickets {
		ticketResponse := MapTicketToTicketResponse(ticket)
		isTicketSaleEnded := IsTicketSaleEndedFor(ticket)
		ticketResponse.IsTicketSaleEnded = isTicketSaleEnded
		ticketResponses = append(ticketResponses, ticketResponse)
	}
	return ticketResponses
}

func MapEventToEventResponse(event *models.Event) *response.EventResponse {
	tickets := GetTicketsFrom(event)
	eventResponse := &response.EventResponse{
		ID:                 event.ID,
		Message:            "event created successfully",
		Name:               event.Name,
		Date:               event.EventDate,
		Time:               event.StartTime,
		ContactInformation: event.ContactInformation,
		Description:        event.Description,
		Status:             event.Status,
		SeriesID:           event.SeriesID,
		Venue:              event.Venue,
		AttendeeTerm:       event.AttendeeTerm,
		EventTheme:         event.EventTheme,
		ImageUrl:           event.ImageUrl,
		Reference:          event.Reference,
		CreatedBy:          event.CreatedBy,
		PublicationState:   event.PublicationState,
		Tickets:            tickets,
	}

	if event.Location != nil {
		eventResponse.Location = &models.Location{
			Longitude: event.Location.Longitude,
			Latitude:  event.Location.Latitude,
			Address:   event.Location.Address,
		}
	}
	return eventResponse
}

func MapTicketToTicketResponse(ticket *models.Ticket) *response.TicketResponse {
	ticketResponse := &response.TicketResponse{
		Type:                         ticket.Type,
		Name:                         ticket.Name,
		Capacity:                     ticket.Capacity,
		Stock:                        ticket.Stock,
		Reference:                    ticket.Reference,
		NumberAvailable:              ticket.NumberAvailable,
		Price:                        ticket.Price,
		PurchaseLimit:                ticket.PurchaseLimit,
		DiscountType:                 ticket.DiscountType,
		AvailableDiscountedTickets:   ticket.AvailableDiscountedTickets,
		IsTransferPaymentFeesToGuest: ticket.IsTransferPaymentFeesToGuest,
		AdditionalInformationFields:  ticket.AdditionalInformationFields,
		Colour:                       ticket.Colour,
		TicketPerks:                  ticket.TicketPerks,
		IsTicketSaleEnded:            ticket.IsSoldOutTicket,
	}
	if ticket.ActivePeriod != nil {
		ticketResponse.SaleEndDate = ticket.ActivePeriod.EndDate
		ticketResponse.SalesEndTime = ticket.ActivePeriod.EndTime
		ticketResponse.SalesStartTime = ticket.ActivePeriod.StartTime
		ticketResponse.SalesStartDate = ticket.ActivePeriod.StartDate
	}
	return ticketResponse
}

func MapCreateAttendeeRequestToAttendee(createAttendeeRequest *dtos.CreateAttendeeRequest) *models.Attendee {

	return &models.Attendee{
		FirstName: createAttendeeRequest.FullName,
		User: &models.User{
			Username: createAttendeeRequest.Username,
			Role:     models.ATTENDEE,
		},
	}
}

func MapAttendeeToAttendeeResponse(attendee *models.Attendee) *response.AttendeeResponse {
	return &response.AttendeeResponse{
		Username:    attendee.Username,
		Message:     "User registered successfully",
		FirstName:   attendee.FirstName,
		LastName:    attendee.LastName,
		PhoneNumber: attendee.PhoneNumber,
	}
}

func IsTicketSaleEndedFor(ticket *models.Ticket) bool {
	if ticket.ActivePeriod == nil {
		return false
	}
	ticketEndTime := ticket.ActivePeriod.EndDate + " " + ticket.ActivePeriod.EndTime
	endTime, err := time.Parse("2006-01-02 15:04:05", ticketEndTime)
	if err != nil {
		log.Println("err: ", err)
		return false
	}
	log.Println("true: ", endTime)
	return time.Now().After(endTime)
}
