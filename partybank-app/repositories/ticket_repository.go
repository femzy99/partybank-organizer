package repositories

import (
	"github.com/djfemz/organizer-service/partybank-app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	crudRepository[models.Ticket, uint64]
	FindAllByEventId(id uint64, pageNumber, pageSize int) ([]*models.Ticket, error)
	FindByNameForEvent(eventId uint64, name string) (*models.Ticket, error)
	FindTicketByReference(reference string) (*models.Ticket, error)
	FindAllTicketsByEventId(eventId uint64) ([]*models.Ticket, error)
	DeleteTicketsFor(eventId uint64) error
	DeleteAllNotIn(eventId uint64, tickets []*models.Ticket) error
}

type raveTicketRepository struct {
	*repositoryImpl[models.Ticket, uint64]
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &raveTicketRepository{&repositoryImpl[models.Ticket, uint64]{
		db,
	}}
}

func (raveTicketRepository *raveTicketRepository) FindAllByEventId(id uint64, pageNumber, pageSize int) ([]*models.Ticket, error) {
	offset, pageSize := getPageInfo(pageNumber, pageSize)
	var tickets []*models.Ticket

	err := raveTicketRepository.Db.Where(&models.Ticket{EventID: id}).Offset(offset).Limit(pageSize).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (raveTicketRepository *raveTicketRepository) FindTicketByReference(reference string) (*models.Ticket, error) {
	var ticket = &models.Ticket{}
	err := raveTicketRepository.Db.Preload(clause.Associations).Where(&models.Ticket{Reference: reference}).First(ticket).Error
	if err != nil {
		return nil, err
	}
	return ticket, nil

}

func (raveTicketRepository *raveTicketRepository) FindAllTicketsByEventId(eventId uint64) ([]*models.Ticket, error) {
	return nil, nil
}

func (raveTicketRepository *raveTicketRepository) DeleteAllNotIn(eventId uint64, tickets []*models.Ticket) error {
	idsToKeep := make([]uint64, 0)
	for _, ticket := range tickets {
		idsToKeep = append(idsToKeep, ticket.ID)
	}
	db := raveTicketRepository.Db
	err := db.Where("event_id = ?", eventId).Where("id NOT IN ?", idsToKeep).Delete(&models.Ticket{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (raveTicketRepository *raveTicketRepository) FindByNameForEvent(eventId uint64, name string) (*models.Ticket, error) {
	ticket := models.Ticket{}
	err := raveTicketRepository.Db.Where(&models.Ticket{EventID: eventId, Name: name}).Find(&ticket).Error
	return &ticket, err
}

func (raveTicketRepository *raveTicketRepository) DeleteTicketsFor(eventId uint64) error {
	err := raveTicketRepository.Db.Where("event_id=?", eventId).Delete(&models.Ticket{}).Error
	return err
}
