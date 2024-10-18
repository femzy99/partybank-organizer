package dtos

import (
	"database/sql/driver"
	"errors"
	"strings"
)

const (
	LIMITED   = "LIMITED"
	UNLIMITED = "UNLIMITED"
)

type AuthRequest struct {
	CreateUserRequest
}

type CreateAttendeeRequest struct {
	FullName    string `json:"full_name,omitempty"`
	Username    string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UpdateAttendeeRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type IssueTicketRequest struct {
	TicketId         uint64 `json:"ticket_id"`
	AttendeeUsername string `json:"attendee_username"`
}

type CreateDiscountRequest struct {
	TicketId uint64  `json:"ticket_id" validate:"required"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Count    uint64  `json:"count"`
	Value    string  `json:"value"`
	Price    float64 `json:"price"`
}

type CreateEventRequest struct {
	Name               string                 `json:"name" validate:"required"`
	Visibility         string                 `json:"visibility" validate:"required"`
	Address            string                 `json:"address" validate:"required"`
	Date               string                 `json:"date" validate:"required"`
	StartTime          string                 `json:"start_time" validate:"required"`
	EndTime            string                 `json:"end_time" validate:"required"`
	ContactInformation string                 `json:"contact_information"`
	Description        string                 `json:"description"`
	SeriesId           uint64                 `json:"series_id" validate:"required"`
	OrganizerId        uint64                 `json:"organizer_id" validate:"required"`
	EventTheme         string                 `json:"event_theme"`
	Latitude           string                 `json:"lat"`
	Longitude          string                 `json:"lng"`
	AttendeeTerm       string                 `json:"-"`
	Venue              string                 `json:"venue" validate:"required"`
	ImageUrl           string                 `json:"image_url"`
	Tickets            []*CreateTicketRequest `json:"tickets"`
}

type UpdateEventRequest struct {
	Name               string `json:"name"`
	Location           string `json:"location"`
	Date               string `json:"date"`
	Time               string `json:"time"`
	ContactInformation string `json:"contact_information"`
	Description        string `json:"description"`
	OrganizerId        uint64 `json:"organizer_id" validate:"required"`
	EventTheme         string `json:"event_theme"`
	MapUrl             string `json:"map_url"`
	MapEmbeddedUrl     string `json:"map_embedded_url"`
	AttendeeTerm       string `json:"attendee_term"`
	Venue              string `json:"venue" validate:"required"`
}

type CreateTicketRequest struct {
	Type                         string      `json:"ticket_type"`
	Name                         string      `json:"name"`
	Capacity                     uint64      `json:"capacity"`
	Stock                        string      `json:"stock"`
	Price                        float64     `json:"price"`
	PurchaseLimit                uint64      `json:"purchase_limit"`
	IsTransferPaymentFeesToGuest bool        `json:"is_transfer_payment_fees_to_guest"` //TODO: Default: false
	EventId                      uint64      `json:"event_id" validate:"required"`
	Colour                       string      `json:"colour"`
	SaleEndDate                  string      `json:"ticket_sale_end_date"`
	SalesEndTime                 string      `json:"ticket_sales_end_time"`
	TicketPerks                  TicketPerks `json:"ticket_perks"`
	PriceChangeDate              string      `json:"price_change_date"`
	PriceChangeTime              string      `json:"price_change_time"`
	SalesStartDate               string      `json:"ticket_sale_start_date"`
	SalesStartTime               string      `json:"ticket_sale_start_time"`
}

type TicketPerks []string

func (o *TicketPerks) Scan(src any) error {
	bytes, ok := src.(string)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(bytes, ",")
	return nil
}
func (o TicketPerks) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,email"`
}

type AddEventStaffRequest struct {
	StaffEmails []string `json:"staff_emails"`
	EventId     uint64   `json:"event_id" validate:"required"`
}

type CreateEventStaffRequest struct {
	StaffEmails []string `json:"staff_emails"`
	EventId     uint64   `json:"event_id" validate:"required"`
}

// "reference": "TCKT002",
// "category": "Single Ticket",
// "type": "Paid",
// "name": "Standard",
// "price": "5000",
// "color": "Yellow",
// "stock": "Limited",
// "capacity": 1000,
// "purchaseLimit": 5,
// "perks": "Free food",
// "salesEndDate": "2024-08-28",
// "salesEndTime": "17:00",
// "priceChangeDate": "2024-08-20",
// "priceChangeTime": "12:00"
type Perks TicketPerks

//type TicketType struct {
//	Reference       string `json:"reference"`
//	Reserved        uint64 `json:"reservedSeats"`
//	MaxSeats        uint64 `json:"maxSeats"`
//	Type            string `json:"type" example:"[FREE/PAID]"`
//	PurchaseLimit   uint64 `json:"purchaseLimit"`
//	Name            string `json:"name"`
//	Price           string `json:"price"`
//	Colour          string `json:"color"`
//	Category        uint64 `json:"category"`
//	Stock           string `json:"stock" example:"[LIMITED/UNLIMITED]"`
//	SalesEndDate    string `json:"salesEndDate"`
//	SalesEndTime    string `json:"salesEndTime"`
//	PriceChangeDate string `json:"priceChangeDate"`
//	PriceChangeTime string `json:"priceChangeTime"`
//	Capacity        uint64 `json:"capacity"`
//	Perks           string `json:"perks"`
//}

type TicketType struct {
	Reference       string `json:"reference"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Price           string `json:"price"`
	Color           string `json:"color"`
	Category        string `json:"category"`
	Stock           string `json:"stock"`
	PurchaseLimit   int    `json:"purchaseLimit"`
	SalesEndDate    string `json:"salesEndDate"`
	SalesEndTime    string `json:"salesEndTime"`
	PriceChangeDate string `json:"priceChangeDate"`
	PriceChangeTime string `json:"priceChangeTime"`
	Capacity        int    `json:"capacity"`
	Perks           string `json:"perks"`
}
type NewTicketMessage struct {
	Reference    string        `json:"eventReference"`
	Types        []*TicketType `json:"ticketTypes"`
	Name         string        `json:"eventName"`
	Venue        string        `json:"venue"`
	TimeFrame    string        `json:"timeFrame"`
	AttendeeTerm string        `json:"attendeeTerm"`
	Date         string        `json:"eventDate"`
	Capacity     uint64        `json:"capacity"`
	CreatedBy    string        `json:"createdBy"`
}

type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Recipient struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailNotificationRequest struct {
	Sender     Sender      `json:"sender"`
	Recipients []Recipient `json:"to"`
	Subject    string      `json:"subject"`
	Content    string      `json:"htmlContent"`
}

type CreateSeriesRequest struct {
	Name        string `json:"name"`
	OrganizerID uint64 `json:"organizer_id" validate:"required"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	SeriesLogo  string `json:"series_logo"`
}

type AttendeeAuthRequest struct {
	Username string `json:"email"`
}

func NewEmailNotificationRequest(recipient, content string) *EmailNotificationRequest {
	return &EmailNotificationRequest{
		Sender: Sender{
			Email: "partybankexperience@gmail.com",
			Name:  "Partybank",
		},
		Recipients: []Recipient{
			{Email: recipient, Name: "Friend"},
		},
		Subject: "Welcome Mail",
		Content: content,
	}
}
