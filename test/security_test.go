package test

//
//import (
//	"github.com/djfemz/rave/partybank-app/models"
//	"github.com/djfemz/rave/partybank-app/security"
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//var organization = &models.Organizer{
//	User: &models.User{
//		Username: "jon@email.com",
//		Role:     models.ORGANIZER,
//	},
//	CreatedAt: time.Now(),
//}
//
//func TestGenerateToken(t *testing.T) {
//	accessToken, err := security.GenerateAccessTokenFor(organization.User)
//	assert.NotNil(t, accessToken)
//	assert.Nil(t, err)
//	assert.NotEmpty(t, accessToken)
//}
//
//func TestExtractUserFromToken(t *testing.T) {
//	//username := "sikiwa1055@glaslack.com"
//	//organizer := &models.Organizer{User: &models.User{Username: username}}
//	//token, _ := security.GenerateAccessTokenFor(organizer.User)
//	//user, _ := security.ExtractUserFrom(token)
//	//log.Println(user)
//	//assert.NotNil(t, user)
//	//assert.NotEmpty(t, user.Username)
//	//assert.Equal(t, user.Username, username)
//}
