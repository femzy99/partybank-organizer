package services

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"os"
)

var (
	clientID          string
	clientSecret      string
	clientCallbackURL string
)

type OauthService struct {
}

func init() {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	clientCallbackURL = os.Getenv("CLIENT_CALLBACK_URL")
}

func NewOauthService() *OauthService {
	goth.UseProviders(google.New(clientID, clientSecret, clientCallbackURL))
	return &OauthService{}
}

func (oauthService *OauthService) authenticate(ctx *gin.Context) {
}
