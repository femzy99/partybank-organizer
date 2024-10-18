package security

import (
	"errors"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"

	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"time"
)

type payload struct {
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

const APP_NAME = "Partybank"

func GenerateAccessTokenFor(user *models.Attendee) (string, error) {
	log.Println("user: ", user)
	token := *jwt.NewWithClaims(jwt.SigningMethodHS256, buildJwtClaimsForAttendee(user))

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateAccessTokenForOrganizer(user *models.User) (string, error) {
	log.Println("user: ", user)
	token := *jwt.NewWithClaims(jwt.SigningMethodHS256, buildJwtClaimsFor(user))

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

//func GenerateAccessTokenForAttendee(user *models.Attendee) (string, error) {
//	log.Println("user: ", user)
//	token := *jwt.NewWithClaims(jwt.SigningMethodHS256, buildJwtClaimsFor(user.User))
//
//	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
//	if err != nil {
//		return "", err
//	}
//	return accessToken, nil
//}

func ExtractUserFrom(token string) (*models.User, error) {
	db := repositories.Connect()
	organizerRepository := repositories.NewOrganizerRepository(db)
	attendeeRepository := repositories.NewAttendeeRepository(db)
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, errors.New("access token is not valid")
	}

	subject, err := tok.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	org, err := organizerRepository.FindByUsername(subject)
	if err != nil {
		att, err := attendeeRepository.FindByUsername(subject)
		if err != nil {
			return nil, err
		}
		return att.User, nil
	}

	return org.User, nil
}

func buildJwtClaimsFor(user *models.User) *jwt.RegisteredClaims {
	return &jwt.RegisteredClaims{
		Issuer:    APP_NAME,
		Subject:   user.Username,
		Audience:  []string{user.Role, strconv.FormatUint(user.ID, 10)},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}
}

func buildJwtClaimsForAttendee(user *models.Attendee) jwt.Claims {
	var claims jwt.MapClaims = make(map[string]interface{})
	claims["firstName"] = user.FirstName
	claims["lastName"] = user.LastName
	claims["phoneNumber"] = user.PhoneNumber
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["iss"] = APP_NAME
	claims["issuedAt"] = jwt.NewNumericDate(time.Now())
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	claims["sub"] = user.Username
	return claims
	//return &jwt.RegisteredClaims{
	//	Issuer:    APP_NAME,
	//	Subject:   user.Username,
	//	Audience:  []string{user.Role, strconv.FormatUint(user.ID, 10)},
	//	IssuedAt:  jwt.NewNumericDate(time.Now()),
	//	ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
	//}
}
