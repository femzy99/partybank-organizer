package dtos

import (
	dtos "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"time"
)

const (
	MAIL_SENDING_SUCCESS_MESSAGE = "mail sent successfully"
	USER_CREATED_SUCCESSFULLY    = "user created successfully"
)

type RaveResponse[T any] struct {
	Data T `json:"data" swaggerignore:"true"`
}

type GoogleUserResponse struct {
	Id              string `json:"id"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"verified_email"`
	Name            string `json:"name"`
}

type CreateDiscountResponse struct {
	Id    uint64  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Code  string  `json:"code,omitempty"`
	Count uint64  `json:"count,omitempty"`
	Value string  `json:"value,omitempty"`
	Price float64 `json:"price,omitempty"`
}

type OrganizationResponse struct {
	*UserResponse
	Name      string            `json:"name,omitempty"`
	CreatedAt time.Time         `json:"created_at,omitempty"`
	Series    []*SeriesResponse `json:"series"`
}

type AttendeeResponse struct {
	Username    string `json:"email"`
	Message     string `json:"message"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type UserResponse struct {
	ID       uint64 `id:"ID" json:"id"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

type SeriesResponse struct {
	ID          uint64           `id:"seriesId" gorm:"primaryKey" json:"series_id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Events      []*EventResponse `json:"events"`
	OrganizerID uint64           `json:"organizer_id,omitempty"`
	ImageUrl    string           `json:"image_url,omitempty"`
	Description string           `json:"description,omitempty"`
	Logo        string           `json:"series_logo"`
}

type LoginResponse struct {
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
	UserID   uint64 `json:"user_id,omitempty"`
}

type CreateOrganizerResponse struct {
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
}

type CreateCalendarResponse struct {
	Message     string           `json:"message,omitempty"`
	ID          uint64           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Events      []*EventResponse `json:"events,omitempty"`
	OrganizerID uint64           `json:"organizer_id,omitempty"`
}

type CalendarResponse struct {
	ID          uint64           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Events      []*EventResponse `json:"events,omitempty"`
	OrganizerID uint64           `json:"organizer_id,omitempty"`
}

type EventResponse struct {
	ID                 uint64            `json:"id"`
	SeriesID           uint64            `json:"series_id,omitempty"`
	SeriesLogo         string            `json:"series_logo,omitempty"`
	Message            string            `json:"message,omitempty"`
	Name               string            `json:"event_name,omitempty"`
	Location           *models.Location  `json:"location,omitempty"`
	Date               string            `json:"date,omitempty"`
	Time               string            `json:"time,omitempty"`
	ContactInformation string            `json:"contact_information,omitempty"`
	Description        string            `json:"description,omitempty"`
	Status             string            `json:"status,omitempty"`
	EventTheme         string            `json:"event_theme,omitempty"`
	MapUrl             string            `json:"map_url,omitempty"`
	MapEmbeddedUrl     string            `json:"map_embedded_url,omitempty"`
	AttendeeTerm       string            `json:"attendee_term,omitempty"`
	Venue              string            `json:"venue,omitempty"`
	ImageUrl           string            `json:"image_url,omitempty"`
	Reference          string            `json:"event_reference,omitempty"`
	CreatedBy          string            `json:"created_by,omitempty"`
	Tickets            []*TicketResponse `json:"tickets"`
	PublicationState   string            `json:"publication_state,omitempty"`
}

type TicketResponse struct {
	Type                         string                             `json:"ticket_type,omitempty"`
	Name                         string                             `json:"name,omitempty"`
	Capacity                     uint64                             `json:"capacity,omitempty"`
	Stock                        string                             `json:"stock"`
	NumberAvailable              uint64                             `json:"number_in_stock,omitempty"`
	Price                        float64                            `json:"price,omitempty"`
	PurchaseLimit                uint64                             `json:"purchase_limit,omitempty"`
	DiscountType                 string                             `json:"discount_type,omitempty"`
	Percentage                   float64                            `json:"percentage,omitempty"`
	DiscountAmount               float64                            `json:"discount_price,omitempty"`
	DiscountCode                 string                             `json:"discount_code,omitempty"`
	AvailableDiscountedTickets   uint64                             `json:"available_discounted_tickets,omitempty"`
	IsTransferPaymentFeesToGuest bool                               `json:"is_transfer_payment_fees_to_guest,omitempty"`
	AdditionalInformationFields  models.AdditionalInformationFields `json:"additional_information_fields,omitempty"`
	Reference                    string                             `json:"ticket_reference,omitempty"`
	Colour                       string                             `json:"colour,omitempty"`
	IsTicketSaleEnded            bool                               `json:"is_ticket_sale_date_expired,omitempty"`
	SaleEndDate                  string                             `json:"ticket_sale_end_date,omitempty"`
	SalesEndTime                 string                             `json:"ticket_sales_end_time,omitempty"`
	SalesStartDate               string                             `json:"ticket_sale_start_date"`
	SalesStartTime               string                             `json:"ticket_sale_start_time"`

	TicketPerks dtos.TicketPerks `json:"ticket_perks,omitempty"`
}
