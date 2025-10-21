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
	appointment := api.Group("/appointment")
	appointment.Get("/swagger/*", swagger.HandlerDefault)

	v1 := appointment.Group("/v1")

	v1.Use(middleware.JwtMiddleware(jwtSvc))

	v1.Get("/patient/history/latest", appointmentHandler.GetLatestAppointmentHistory)
	v1.Get("/patient/history", appointmentHandler.GetPatientAppointmentHistory)
	v1.Get("/patient/incoming", appointmentHandler.IncomingAppointmentOfPatient)
	v1.Post("/patient", appointmentHandler.BookAppointment)

	v1.Get("/doctor/:doctor_id/slots", appointmentHandler.GetDoctorSlots)

	v1.Post("/doctor/shift", appointmentHandler.CreateDoctorShift)
	v1.Delete("/doctor/shift", appointmentHandler.DeleteDoctorShift)
	v1.Get("/doctor/shift", appointmentHandler.GetDoctorActiveShifts)
	v1.Get("/doctor", appointmentHandler.GetDoctorIncomingAppointments)

	v1.Post("/cancel", appointmentHandler.CancelAppointment)
}
