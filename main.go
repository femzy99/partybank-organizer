package main

import (
	_ "github.com/djfemz/organizer-service/docs"
	"github.com/djfemz/organizer-service/partybank-app/config"
	handlers "github.com/djfemz/organizer-service/partybank-app/controllers"
	"github.com/djfemz/organizer-service/partybank-app/repositories"
	"github.com/djfemz/organizer-service/partybank-app/security/controllers"
	"github.com/djfemz/organizer-service/partybank-app/security/middlewares"
	services2 "github.com/djfemz/organizer-service/partybank-app/security/services"
	"github.com/djfemz/organizer-service/partybank-app/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
)

var err error
var db *gorm.DB
var eventRepository repositories.EventRepository
var organizerRepository repositories.OrganizerRepository
var seriesRepository repositories.SeriesRepository
var ticketRepository repositories.TicketRepository
var eventStaffRepository repositories.EventStaffRepository
var attendeeRepository repositories.AttendeeRepository

var eventService services.EventService
var organizerService services.OrganizerService
var seriesService services.SeriesService
var ticketService services.TicketService
var eventStaffService services.EventStaffService
var authService *services2.AuthService
var attendeeService services.AttendeeService

var organizerController *handlers.OrganizerController
var eventController *handlers.EventController
var seriesController *handlers.SeriesController
var ticketController *handlers.TicketController
var attendeeController *handlers.AttendeeController
var authController *controllers.AuthController

var objectValidator *validator.Validate

func init() {
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading configuration: ", err)
	}
	log.Println("connnecting to db")
	db = repositories.Connect()
}

//partybank-organizer-269c8057a65f.herokuapp.com

// @title           Partybank Organizer Service
// @version         1.0
// @description     Partybank Organizer Service.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    https://www.thepartybank.com
// @contact.email  unavailable
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host partybank-organizer-269c8057a65f.herokuapp.com
// @schemes https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
func main() {
	config.GoogleConfig()
	router := gin.Default()
	configureAppComponents()
	middlewares.Routers(router, organizerController,
		eventController, seriesController, ticketController,
		authService, attendeeController, authController, attendeeRepository)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	err = router.Run(":" + port)
	if err != nil {
		log.Println("Error starting server: ", err)
	}
}

func configureAppComponents() {
	objectValidator = validator.New()
	configureRepositoryComponents()
	configureServiceComponents()
	configureControllers()
}

func configureControllers() {
	organizerController = handlers.NewOrganizerController(organizerService, objectValidator)
	eventController = handlers.NewEventController(eventService, objectValidator)
	seriesController = handlers.NewSeriesController(seriesService, objectValidator)
	ticketController = handlers.NewTicketController(ticketService, objectValidator)
	attendeeController = handlers.NewAttendeeController(attendeeService, objectValidator)
	authController = controllers.NewAuthController(authService)
}

func configureServiceComponents() {
	mailService := services.NewMailService()
	seriesService = services.NewSeriesService(seriesRepository)
	eventStaffService = services.NewEventStaffService(eventStaffRepository, eventRepository)
	organizerService = services.NewOrganizerService(organizerRepository, eventStaffService, seriesService, ticketService, attendeeService)
	eventService = services.NewEventService(eventRepository, organizerService, seriesService, ticketService)
	ticketService = services.NewTicketService(ticketRepository, eventService)
	eventService.SetTicketService(ticketService)
	attendeeService = services.NewAttendeeService(attendeeRepository, mailService)
	authService = services2.NewAuthService(organizerService, attendeeService, mailService)
}

func configureRepositoryComponents() {
	eventRepository = repositories.NewEventRepository(db)
	organizerRepository = repositories.NewOrganizerRepository(db)
	seriesRepository = repositories.NewSeriesRepository(db)
	ticketRepository = repositories.NewTicketRepository(db)
	eventStaffRepository = repositories.NewEventStaffRepository(db)
	attendeeRepository = repositories.NewAttendeeRepository(db)
}
