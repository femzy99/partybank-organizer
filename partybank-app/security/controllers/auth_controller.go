package controllers

import (
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	authService "github.com/djfemz/organizer-service/partybank-app/security/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var err error

type AuthController struct {
	*authService.AuthService
}

func NewAuthController(authService *authService.AuthService) *AuthController {
	return &AuthController{
		authService,
	}
}

// AuthHandler godoc
// @Summary      Authenticate user
// @Description  Authenticate user
// @Tags         Auth
// @Accept       json
// @Param 		 tags body dtos.AuthRequest true "Auth tags"
// @Produce      json
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Failure      500  {object}  dtos.RaveResponse
// @Router       /auth/login [post]
func (authController *AuthController) AuthHandler(ctx *gin.Context) {
	var signInRequest request.AuthRequest
	if err = ctx.BindJSON(&signInRequest); err != nil {
		handleError(ctx, err)
		return
	}
	res, err := authController.AuthService.Authenticate(&signInRequest)
	if res != nil {
		ctx.JSON(http.StatusOK, response.RaveResponse[response.LoginResponse]{Data: *res})
	} else {
		handleError(ctx, err)
	}
}

// ValidateOtp godoc
// @Summary      Validate Otp
// @Description  Validate Otp
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        code   query   int  true  "otp code"
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Failure      500  {object}  dtos.RaveResponse
// @Router       /auth/otp/validate [get]
func (authController *AuthController) ValidateOtp(ctx *gin.Context) {
	code := ctx.Query("code")
	res, err := authController.AuthService.ValidateOtp(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.LoginResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// AuthenticateAttendee godoc
// @Summary      Authenticate attendee
// @Description  Authenticate attendee
// @Tags         Auth
// @Accept       json
// @Param 		 tags body dtos.AttendeeAuthRequest true "Auth tags"
// @Produce      json
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Failure      500  {object}  dtos.RaveResponse
// @Router       /auth/login/attendee [post]
func (authController *AuthController) AuthenticateAttendee(ctx *gin.Context) {
	var signInRequest = request.AttendeeAuthRequest{}
	if err = ctx.BindJSON(&signInRequest); err != nil {
		handleError(ctx, err)
		return
	}
	authResponse, err := authController.AuthService.AuthenticateAttendee(signInRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.RaveResponse[*response.LoginResponse]{Data: authResponse})
}

func handleError(ctx *gin.Context, err error) {
	ctx.IndentedJSON(http.StatusBadRequest, &response.RaveResponse[string]{Data: err.Error()})
}
