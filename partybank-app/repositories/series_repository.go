package repositories

import (
	"github.com/djfemz/organizer-service/partybank-app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeriesRepository interface {
	Repository[models.Series, uint64]
	FindPublicSeriesFor(organizer uint64) (*models.Series, error)
	FindAllSeriesFor(organizer uint64, pageNumber int, pageSize int) ([]*models.Series, error)
}

type raveCalendarRepository struct {
	repositoryImpl[models.Series, uint64]
}

func NewSeriesRepository(db *gorm.DB) SeriesRepository {
	return &raveCalendarRepository{
		repositoryImpl[models.Series, uint64]{
			db,
		},
	}
}

func (raveCalendarRepository *raveCalendarRepository) FindPublicSeriesFor(organizer uint64) (*models.Series, error) {
	foundSeries := &models.Series{}
	err := raveCalendarRepository.Db.Where(&models.Series{Name: "Public", OrganizerID: organizer}).Find(foundSeries).Error
	if err != nil {
		return nil, err
	}
	return foundSeries, nil
}

func (raveCalendarRepository *raveCalendarRepository) FindAllSeriesFor(organizer uint64, pageNumber int, pageSize int) ([]*models.Series, error) {
	if pageSize < 1 {
		pageSize = 1
	}
	if pageNumber < 1 {
		pageNumber = 1
	} else if pageSize > 100 {
		pageSize = 100
	}
	offset := (pageNumber - 1) * pageSize
	userSeries := make([]*models.Series, 0)
	err := raveCalendarRepository.Db.Preload(clause.Associations).Where(&models.Series{OrganizerID: organizer}).Offset(offset).Limit(pageSize).Find(&userSeries).Error
	if err != nil {
		return nil, err
	}
	return userSeries, nil
}
