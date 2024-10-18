package repositories

import (
	"github.com/djfemz/organizer-service/partybank-app/models"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	crudRepository[models.Discount, uint64]
}

type discountRepository struct {
	*repositoryImpl[models.Discount, uint64]
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepository{
		&repositoryImpl[models.Discount, uint64]{
			db,
		},
	}
}
