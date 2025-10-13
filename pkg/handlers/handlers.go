package handlers

import (
	"appointment-service/pkg/dto"
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

func (h *AppointmentHandler) BookAppointment(c *fiber.Ctx) error {
	patientID := c.Locals("userID").(string)
	var body dto.BookAppointmentRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	resp, err := h.appointmentService.BookAppointment(patientID, &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
