package repositories

import "github.com/djfemz/organizer-service/partybank-app/models"

type IssuedTicketRepository interface {
	crudRepository[models.IssuedTicket, uint64]
}

type IssuedTicketRepositoryImpl struct {
}

func NewIssuedTicketRepository() IssuedTicketRepository {
	return &repositoryImpl[models.IssuedTicket, uint64]{}
}
