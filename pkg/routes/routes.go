package routes

import (
	// "appointment-service/pkg/context"
	// "appointment-service/pkg/dto"
	_ "appointment-service/docs"
	"appointment-service/pkg/handlers"
	"appointment-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, appointmentHandler *handlers.AppointmentHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")
	api.Get("/swagger/*", swagger.HandlerDefault)
	appointment := api.Group("/appointment")
	v1 := appointment.Group("/v1")


}
