package controllers

import (
	"encoding/json"
	"github.com/djfemz/organizer-service/partybank-app/config"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/models"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/security"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var (
	clientState string
)

const tokenEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type OauthController struct {
	repositories.AttendeeRepository
}

func init() {
	clientState = os.Getenv("GOOGLE_CLIENT_STATE")
}

// GoogleLogin godoc
// @Summary      Authenticate with Google
// @Description  Sign in with Google
// @Tags         Google
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Router       /auth/google/login [get]
func (oauthController *OauthController) GoogleLogin(ctx *gin.Context) {
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(clientState)
	log.Println("url: ", url)
	ctx.JSON(http.StatusOK, response.RaveResponse[string]{Data: url})
}

func (oauthController *OauthController) GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	if state != clientState {
		ctx.JSON(http.StatusBadRequest, "States don't Match!!")
		return
	}
	code := ctx.Query("code")
	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(ctx, code)
	if err != nil {
		log.Println("Error: ", err)
		ctx.JSON(http.StatusBadRequest, "Code-Token Exchange Failed")
		return
	}
	log.Println("tokenEndpoint: ", tokenEndpoint)
	resp, err := http.Get(tokenEndpoint + token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "User Data Fetch Failed")
		return
	}

	var googleUser = response.GoogleUserResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&googleUser)
	if err != nil {
		log.Println("error extracting user data", err.Error())
		ctx.JSON(http.StatusBadRequest, "error extracting user data")
		return
	}
	user := &models.User{Username: googleUser.Email, Role: models.USER}
	attendee, err := oauthController.AttendeeRepository.FindByUsername(googleUser.Email)
	var accessToken string
	if err != nil || attendee == nil {
		log.Println("error: ", err)
		attendee = &models.Attendee{User: user}
	}
	accessToken, err = security.GenerateAccessTokenFor(attendee)
	if err != nil {
		log.Println("error: ", err)
		ctx.JSON(http.StatusBadRequest, "failed to generate access token")
		return
	}
	log.Println("token: ", accessToken)
	data := utils.FRONT_END_PROD_URL + "/validate-token?token=" + accessToken + "&type=google"
	log.Println("uri: ", data)
	ctx.Redirect(307, data)
}
