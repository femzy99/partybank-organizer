package repositories

import (
	"github.com/djfemz/organizer-service/partybank-app/models"
	"gorm.io/gorm"
)

type AttendeeRepository interface {
	crudRepository[models.Attendee, uint64]
	FindByUsername(username string) (*models.Attendee, error)
}

type raveAttendeeRepository struct {
	*repositoryImpl[models.Attendee, uint64]
}

func NewAttendeeRepository(db *gorm.DB) AttendeeRepository {
	return &raveAttendeeRepository{
		&repositoryImpl[models.Attendee, uint64]{
			db,
		},
	}
}

func (attendeeRepository *raveAttendeeRepository) FindByUsername(username string) (*models.Attendee, error) {
	attendee := &models.Attendee{}
	err := attendeeRepository.Db.Where("username=?", username).First(attendee).Error
	if err != nil {
		return nil, err
	}
	return attendee, err
}
