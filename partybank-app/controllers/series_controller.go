package controllers

import (
	request "github.com/djfemz/organizer-service/partybank-app/dtos/request"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/services"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type SeriesController struct {
	services.SeriesService
	*validator.Validate
}

func NewSeriesController(seriesService services.SeriesService, objectValidator *validator.Validate) *SeriesController {
	return &SeriesController{
		seriesService,
		objectValidator,
	}
}

// CreateSeries godoc
// @Summary      Add New Series
// @Description  Adds New Series
// @Tags         Series
// @Accept       json
// @Param 		 tags body dtos.CreateSeriesRequest true "Series tags"
// @Produce      json
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/series [post]
func (seriesController *SeriesController) CreateSeries(ctx *gin.Context) {
	createSeriesRequest := &request.CreateSeriesRequest{}
	err := ctx.BindJSON(createSeriesRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}
	err = seriesController.Struct(createSeriesRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}
	resp, err := seriesController.AddSeries(createSeriesRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}
	ctx.JSON(http.StatusCreated, &response.RaveResponse[*response.CreateCalendarResponse]{Data: resp})
}

// GetSeriesById godoc
// @Summary      Get Series by id
// @Description  Get Series by id
// @Tags         Series
// @Accept       json
// @Produce      json
// @Param        id   path   int  true  "series id"
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Failure      500  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/series/{id} [get]
func (seriesController *SeriesController) GetSeriesById(ctx *gin.Context) {
	seriesService := seriesController.SeriesService
	id, err := extractParamFromRequest("id", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}
	calendar, err := seriesService.GetCalendar(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.RaveResponse[error]{Data: err})
		return
	}

	ctx.JSON(http.StatusOK, &response.RaveResponse[*response.SeriesResponse]{Data: calendar})

}

// GetSeriesForOrganizer godoc
// @Summary      Get Series by organizerId
// @Description  Get Series by organizerId
// @Tags         Series
// @Accept       json
// @Produce      json
// @Param        organizerId   path   int  true  "organizerId"
// @Param        page   query   int  true  "page"
// @Param        size   query   int  true  "size"
// @Success      200  {object}  dtos.RaveResponse
// @Failure      400  {object}  dtos.RaveResponse
// @Security Bearer
// @Router       /api/v1/series/organizer/{organizerId} [get]
func (seriesController *SeriesController) GetSeriesForOrganizer(ctx *gin.Context) {
	organizerId, err := extractParamFromRequest("organizerId", ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}
	log.Println("org id", organizerId)
	page := ctx.Query("page")
	size := ctx.Query("size")
	pageNumber, err := utils.ConvertQueryStringToInt(page)
	if err != nil {
		handleError(ctx, err)
		return
	}
	pageSize, err := utils.ConvertQueryStringToInt(size)
	if err != nil {
		handleError(ctx, err)
		return
	}
	orgSeries, err := seriesController.SeriesService.GetSeriesFor(organizerId, pageNumber, pageSize)
	if err != nil {
		handleError(ctx, err)
		return
	}
	log.Println("series for org: ", orgSeries)
	ctx.JSON(http.StatusOK, &response.RaveResponse[[]*response.SeriesResponse]{Data: orgSeries})
}
