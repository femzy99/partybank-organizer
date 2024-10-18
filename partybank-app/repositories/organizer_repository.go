package repositories

import (
	"github.com/djfemz/organizer-service/partybank-app/models"
	otp2 "github.com/djfemz/organizer-service/partybank-app/security/otp"
	"gorm.io/gorm"
	"log"
)

type OrganizerRepository interface {
	crudRepository[models.Organizer, uint64]
	FindByUsername(username string) (*models.Organizer, error)
	FindByOtp(otp string) (*models.Organizer, error)
}

type organizerRepositoryImpl struct {
	repositoryImpl[models.Organizer, uint64]
}

func NewOrganizerRepository(db *gorm.DB) OrganizerRepository {
	var organizerRepository OrganizerRepository = &organizerRepositoryImpl{
		repositoryImpl[models.Organizer, uint64]{
			db,
		},
	}
	return organizerRepository
}

func (organizerRepository *organizerRepositoryImpl) FindByUsername(username string) (*models.Organizer, error) {

	var organization = new(models.Organizer)
	err := organizerRepository.Db.Where("username=?", username).First(&organization).Error
	return organization, err
}

func (organizerRepository *organizerRepositoryImpl) FindByOtp(otp string) (*models.Organizer, error) {
	var organizer models.Organizer

	err := organizerRepository.Db.Where(&models.Organizer{Otp: &otp2.OneTimePassword{Code: otp}}).Find(&organizer).Error
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}
	log.Println("found org: ", organizer)
	return &organizer, nil
}
