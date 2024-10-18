package test

//
//import (
//	dtos "github.com/djfemz/rave/partybank-app/dtos/request"
//	"github.com/djfemz/rave/partybank-app/services"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//var discountService services.DiscountService = services.NewDiscountService()
//
//func TestCreateDiscount(t *testing.T) {
//	createDiscountRequest := &dtos.CreateDiscountRequest{
//		TicketId: 1,
//		Name:     "test discount",
//		Code:     "ABCD",
//		Count:    30,
//		Value:    "3%",
//		Price:    2000.00,
//	}
//	response, err := discountService.CreateDiscount(createDiscountRequest)
//	assert.NotNil(t, response)
//	assert.Nil(t, err)
//}
