package middlewares

import (
	handlers "github.com/djfemz/organizer-service/partybank-app/controllers"
	response "github.com/djfemz/organizer-service/partybank-app/dtos/response"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/security"
	"github.com/djfemz/organizer-service/partybank-app/security/controllers"
	"github.com/djfemz/organizer-service/partybank-app/security/services"
	"github.com/djfemz/organizer-service/partybank-app/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"strings"
	"time"
)

var routesAuthorities map[string][]string

func Routers(router *gin.Engine, organizerController *handlers.OrganizerController,
	eventController *handlers.EventController, seriesController *handlers.SeriesController,
	ticketController *handlers.TicketController, authService *services.AuthService,
	attendeeController *handlers.AttendeeController, authController *controllers.AuthController, attendeeRepo repositories.AttendeeRepository) {

	router.Use(cors.New(configureCors()))
	protected := router.Group("/api/v1", AuthMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	{
		protected.POST("/event", eventController.CreateEvent)
		protected.GET("/event/:id", eventController.GetEventById)
		protected.PUT("/event/:id", eventController.EditEvent)
		protected.GET("/event/organizer", eventController.GetAllEventsForOrganizer)
		protected.POST("/event/staff", organizerController.AddEventStaff)
		protected.POST("/ticket", ticketController.AddTicketToEvent)
		protected.GET("/ticket/:eventId", ticketController.GetAllTicketsForEvent)
		protected.GET("/ticket", ticketController.GetTicketById)
		protected.POST("/series", seriesController.CreateSeries)
		protected.GET("/series/:id", seriesController.GetSeriesById)
		protected.GET("/series/organizer/:organizerId", seriesController.GetSeriesForOrganizer)
		protected.GET("/event/publish/:id", eventController.PublishEvent)
		protected.PUT("/attendee/update/:username", attendeeController.UpdateAttendee)
	}

	oauthController := &controllers.OauthController{
		AttendeeRepository: attendeeRepo,
	}
	router.POST("/auth/login", authController.AuthHandler)
	router.POST("/auth/login/attendee", authController.AuthenticateAttendee)
	router.GET("/auth/google/login", oauthController.GoogleLogin)
	router.GET("/auth/google/redirect", oauthController.GoogleCallback)
	router.GET("/auth/otp/validate", authController.ValidateOtp)
	router.GET("/api/v1/event/discover", eventController.DiscoverEvents)
	router.GET("/api/v1/event/reference/:reference", eventController.GetEventByReference)
	router.GET("/api/v1/ticket/update", ticketController.UpdateTicketSoldOutStatusByReference)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(utils.AUTHORIZATION)
		log.Println("auth header: ", authHeader)
		token := extractTokenFrom(authHeader)
		if !isValid(token) {
			ctx.AbortWithStatusJSON(http.StatusForbidden,
				&response.RaveResponse[string]{Data: "access token is invalid"})
			return
		}
		ctx.Next()
	}
}

func isValid(token string) bool {
	org, err := security.ExtractUserFrom(token)
	log.Println("\norg: ", org, "\nerr: ", err)
	if err != nil || org == nil {
		return false
	}
	return true
}

func extractTokenFrom(authHeader string) string {
	authValue := strings.Split(authHeader, " ")
	token := authValue[len(authValue)-1]
	log.Println("token: ", token)
	return token
}

func configureCors() cors.Config {
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"Content-Length"},
		AllowHeaders:     []string{utils.AUTHORIZATION, "Content-Type", "Origin", "Accept"},
		AllowMethods:     []string{http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodGet},
	}
	return config
}
