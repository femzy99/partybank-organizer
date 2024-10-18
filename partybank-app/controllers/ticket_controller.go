package controllers

import (
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/services"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type TicketController struct {
	services.TicketService
	*validator.Validate
}

func NewTicketController(ticketService services.TicketService, objectValidator *validator.Validate) *TicketController {
	return &TicketController{
		ticketService,
		objectValidator,
	}
}

// AddTicketToEvent godoc
// @Summary      Add Ticket to Event
// @Description  Add Ticket to Event
// @Tags         Tickets
// @Accept       json
// @Param 		 tags body dtos.CreateTicketRequest true "Ticket tags"
// @Produce      json
// @Success      200  {object}  dtos.TicketResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/ticket [post]
func (ticketController *TicketController) AddTicketToEvent(ctx *gin.Context) {
	addTicketRequest := &request.CreateTicketRequest{}
	err := ctx.BindJSON(addTicketRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	err = ticketController.Struct(addTicketRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ticketResponse, err := ticketController.CreateTicketFor(addTicketRequest)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, ticketResponse)
}

// GetAllTicketsForEvent godoc
// @Summary      Get all Tickets for Event
// @Description   Get all Tickets for Event
// @Tags         Tickets
// @Accept       json
// @Param        eventId path int  true  "eventId"
// @Param        page   query   int  true  "page"
// @Param        size   query   int  true  "size"
// @Produce      json
// @Success      200  {array}  models.Ticket
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/ticket/{eventId} [get]
func (ticketController *TicketController) GetAllTicketsForEvent(ctx *gin.Context) {
	eventId, err := extractParamFromRequest("eventId", ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	page := ctx.Query("page")
	size := ctx.Query("size")
	pageNumber, err := utils.ConvertQueryStringToInt(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}
	pageSize, err := utils.ConvertQueryStringToInt(size)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}
	tickets, err := ticketController.GetAllTicketsFor(eventId, pageNumber, pageSize)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, tickets)
}

// GetTicketById godoc
// @Summary      Get Ticket By id
// @Description   Get Ticket By id
// @Tags         Tickets
// @Accept       json
// @Param        ticketId  query int  true  "ticketId"
// @Produce      json
// @Success      200  {object}  dtos.TicketResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/ticket [get]
func (ticketController *TicketController) GetTicketById(ctx *gin.Context) {
	eventId, err := strconv.ParseUint(ctx.Query("ticketId"), 10, 64)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ticket, err := ticketController.GetById(eventId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, ticket)
}

// UpdateTicketSoldOutStatusByReference godoc
// @Summary      Update Ticket By reference
// @Description  Update Ticket By reference
// @Tags         Tickets
// @Accept       json
// @Param        reference query string  true  "reference"
// @Produce      json
// @Success      200  {object}  dtos.TicketResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/ticket/update [get]
func (ticketController *TicketController) UpdateTicketSoldOutStatusByReference(ctx *gin.Context) {
	reference := ctx.Query("reference")
	ticket, err := ticketController.TicketService.UpdateTicketSoldOutBy(reference)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, ticket)
}
