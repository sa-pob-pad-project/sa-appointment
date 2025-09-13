package handlers

import "appointment-service/pkg/service"

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}
