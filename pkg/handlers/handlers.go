package handlers

import (
	"appointment-service/pkg/service"

	"github.com/gofiber/fiber/v2"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}


func (h *AppointmentHandler) GetDoctorSlots(c *fiber.Ctx) error {
	doctorID := c.Params("doctor_id")
	slots, err := h.appointmentService.GetDoctorSlots(doctorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(slots)
}
