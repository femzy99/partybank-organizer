package models

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
	"time"

	dtos "github.com/djfemz/organizer-service/partybank-app/dtos/request"

	"gorm.io/gorm"

	"github.com/djfemz/organizer-service/partybank-app/security/otp"
)

var Entities = make(map[string]any, 100)

const (
	ATTENDEE    = "ATTENDEE"
	ORGANIZER   = "ORGANIZER"
	EVENT_STAFF = "EVENT_STAFF"
	USER        = "USER"
)

const (
	PAST      = "PAST"
	UPCOMING  = "UPCOMING"
	LIMITED   = "Limited"
	UNLIMITED = "Unlimited"
	PAID      = "Paid"
)

const (
	PUBLISHED = "PUBLISHED"
	DRAFT     = "DRAFT"
)

const (
	ACTIVE    = "ACTIVE"
	SUSPENDED = "SUSPENDED"
	IN_ACTIVE = "IN_ACTIVE"
)

// Used to register entities
func init() {
	Entities[reflect.ValueOf(Event{}).String()] = Event{}
	Entities[reflect.ValueOf(Organizer{}).String()] = Organizer{}
	Entities[reflect.ValueOf(EventStaff{}).String()] = EventStaff{}
	Entities[reflect.ValueOf(Ticket{}).String()] = Ticket{}
	Entities[reflect.ValueOf(Series{}).String()] = Series{}
}

type Organizer struct {
	ID uint64 `id:"ID" gorm:"primaryKey"`
	*User
	Name      string
	CreatedAt time.Time
	Otp       *otp.OneTimePassword `gorm:"embedded;embeddedPrefix:otp"`
	EventId   uint64
	Series    []*Series
}

type User struct {
	ID       uint64 `id:"ID" gorm:"primaryKey" json:"id"`
	Username string `json:"username" gorm:"unique"`
	Role     string `json:"role"`
}

type Attendee struct {
	ID        uint64 `id:"ID" gorm:"primaryKey" json:"id"`
	FirstName string
	LastName  string
	*User
	PhoneNumber string
	IsActive    bool
}

type IssuedTicket struct {
	ID     uint64 `id:"ID" gorm:"primaryKey" json:"id"`
	Issuer *Organizer
	*gorm.Model
	Attendee *Attendee
	Ticket   *Ticket
}

type AdditionalInformationFields []string

func (o *AdditionalInformationFields) Scan(src any) error {
	bytes, ok := src.(string)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(bytes, ",")
	return nil
}
func (o AdditionalInformationFields) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}

type Ticket struct {
	ID                           uint64 `id:"id" gorm:"primaryKey"`
	Type                         string
	Name                         string                      `json:"name"`
	Capacity                     uint64                      `json:"capacity"`
	Category                     string                      `json:"category"`
	Stock                        string                      `json:"stock"`
	GroupTicketCapacity          uint64                      `json:"group_ticket_capacity"`
	NumberAvailable              uint64                      `json:"number_available"`
	Price                        float64                     `json:"price"`
	PurchaseLimit                uint64                      `json:"purchase_limit"`
	DiscountType                 string                      `json:"discount_type,-"`
	Percentage                   float64                     `json:"percentage,-"`
	DiscountAmount               float64                     `json:"discount_price,-"`
	DiscountCode                 string                      `json:"discount_code,-"`
	AvailableDiscountedTickets   uint64                      `json:"available_discounted_tickets,-"`
	AdditionalInformationFields  AdditionalInformationFields `gorm:"type:VARCHAR(255)" json:"additional_information_fields,omitempty"`
	TicketPerks                  dtos.TicketPerks            `gorm:"type:VARCHAR(255)" json:"ticket_perks"`
	IsTransferPaymentFeesToGuest bool                        `json:"is_transfer_payment_fees_to_guest"`
	EventID                      uint64                      `json:"event_id"`
	Reference                    string                      `json:"reference"`
	Colour                       string                      `json:"colour"`
	ActivePeriod                 *ActivePeriod               `gorm:"embedded" json:"active_period"`
	IsSoldOutTicket              bool                        `json:"is_sold_out_ticket"`
	MaxSeats                     uint64                      `json:"max_seats"`
	EventReference               string                      `json:"event_reference"`
	Reserved                     uint64                      `json:"reserved"`
	DeletedAt                    gorm.DeletedAt
}

type ActivePeriod struct {
	StartDate string `json:"ticket_sale_start_date"`
	StartTime string `json:"ticket_sale_start_time"`
	EndDate   string `json:"ticket_sale_end_date"`
	EndTime   string `json:"ticket_sale_end_time"`
}

type Event struct {
	ID                    uint64    `id:"id" gorm:"primaryKey" json:"id"`
	Name                  string    `json:"name"`
	Location              *Location `json:"location" gorm:"embedded"`
	DeletedAt             gorm.DeletedAt
	EventDate             string `json:"date"`
	StartTime             string `json:"event_start"`
	EndTime               string `json:"event_end"`
	ContactInformation    string `json:"contact_information"`
	ImageUrl              string `json:"image_url"`
	Description           string `json:"description"`
	SeriesID              uint64 `json:"series_id"`
	Status                string `json:"status"`
	EventStaffID          uint64 `json:"event_staff_id"`
	IsNotificationEnabled bool   `json:"is_notification_enabled"`
	TicketID              uint64 `json:"ticket_id"`
	EventTheme            string `json:"event_theme"`
	AttendeeTerm          string `json:"-"`
	Venue                 string `json:"venue"`
	Reference             string `json:"event_reference"`
	Tickets               []*Ticket
	EventStaff            []*EventStaff
	CreatedBy             string
	PublicationState      string
	DateCreated           string `json:"created_at"`
	IsEventDeleted        bool
}

type Location struct {
	Longitude string `json:"lng,omitempty"`
	Latitude  string `json:"lat,omitempty"`
	Address   string `json:"address,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
}

type Series struct {
	ID   uint64 `id:"id" gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	gorm.Model
	Events      []*Event `json:"events"`
	OrganizerID uint64   `json:"organizer_id"`
	ImageUrl    string   `json:"image_url"`
	Description string   `json:"description"`
	Logo        string   `json:"series_logo"`
}

type EventStaff struct {
	ID      uint64 `id:"id" gorm:"primaryKey" json:"id"`
	*User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	EventID uint64
}

type Discount struct {
	ID uint64 `id:"id" gorm:"primaryKey" json:"id"`
	*Ticket
	Name  string
	Code  string
	Count uint64
	Value string
	Price float64
}
