package controllers

import (
	dtos "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AttendeeController struct {
	attendeeService services.AttendeeService
	objectValidator *validator.Validate
}

func NewAttendeeController(attendeeService services.AttendeeService, objectValidator *validator.Validate) *AttendeeController {
	return &AttendeeController{
		attendeeService: attendeeService,
		objectValidator: objectValidator,
	}
}

//// Register godoc
//// @Summary      Onboard Attendee
//// @Description  Onboard Attendee
//// @Tags         Auth
//// @Accept       json
//// @Param 		 tags body dtos.CreateAttendeeRequest true "Auth tags"
//// @Produce      json
//// @Success      201  {object}  dtos.RaveResponse
//// @Failure      400  {object}  dtos.RaveResponse
//// @Router       /auth/attendee  [post]
//func (attendeeController AttendeeController) Register(ctx *gin.Context) {
//	attendeeAuthRequest := &dtos.CreateAttendeeRequest{}
//	err := ctx.BindJSON(attendeeAuthRequest)
//	if err != nil {
//		handleError(ctx, err)
//		return
//	}
//	err = attendeeController.objectValidator.Struct(attendeeAuthRequest)
//	if err != nil {
//		handleError(ctx, err)
//		return
//	}
//	res, err := attendeeController.attendeeService.Register(attendeeAuthRequest)
//	if err != nil {
//		handleError(ctx, err)
//		return
//	}
//	ctx.JSON(http.StatusCreated, response.RaveResponse[*response.AttendeeResponse]{Data: res})
//}

// UpdateAttendee godoc
// @Summary      Update Attendee
// @Description  Update Attendee Details
// @Tags         Attendee
// @Accept       json
// @Param 		 tags body dtos.UpdateAttendeeRequest true "Attendee tags"
// @Param 		 username  path string  true  "username"
// @Produce      json
// @Success      201  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/attendee/update/{username}  [put]
func (attendeeController AttendeeController) UpdateAttendee(ctx *gin.Context) {
	username := ctx.Param("username")

	attendeeAuthRequest := &dtos.UpdateAttendeeRequest{}
	err := ctx.BindJSON(attendeeAuthRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	err = attendeeController.objectValidator.Struct(attendeeAuthRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	res, err := attendeeController.attendeeService.UpdateAttendee(username, attendeeAuthRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response.RaveResponse[*response.AttendeeResponse]{Data: res})
}
