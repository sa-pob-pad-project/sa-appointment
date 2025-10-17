package routes

import (
	// "appointment-service/pkg/context"
	// "appointment-service/pkg/dto"
	_ "appointment-service/docs"
	"appointment-service/pkg/handlers"
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, appointmentHandler *handlers.AppointmentHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")
	api.Get("/swagger/*", swagger.HandlerDefault)
	appointment := api.Group("/appointment")
	v1 := appointment.Group("/v1")

	v1.Use(middleware.JwtMiddleware(jwtSvc))

	v1.Get("history", appointmentHandler.GetPatientAppointmentHistory)
	v1.Get("/incoming", appointmentHandler.IncomingAppointmentOfPatient)
	v1.Post("/", appointmentHandler.BookAppointment)

	v1.Get("/:doctor_id/slots", appointmentHandler.GetDoctorSlots)

	v1.Post("/doctor/shift", appointmentHandler.CreateDoctorShift)
	v1.Delete("/doctor/shift", appointmentHandler.DeleteDoctorShift)

}
