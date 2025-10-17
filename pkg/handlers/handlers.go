package handlers

import (
	"appointment-service/pkg/apperr"
	contextUtils "appointment-service/pkg/context"
	"appointment-service/pkg/dto"
	"appointment-service/pkg/response"
	service "appointment-service/pkg/services"

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
func (h *AppointmentHandler) GetPatientAppointmentHistory(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointments, err := h.appointmentService.GetPatientAppointmentHistory(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, appointments)
}

func (h *AppointmentHandler) IncomingAppointmentOfPatient(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointments, err := h.appointmentService.GetPatientIncomingAppointments(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, appointments)
}

func (h *AppointmentHandler) GetDoctorSlots(c *fiber.Ctx) error {
	doctorID := c.Params("doctor_id")
	slots, err := h.appointmentService.GetDoctorSlots(doctorID)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, slots)
}

func (h *AppointmentHandler) BookAppointment(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	var body dto.BookAppointmentRequest
	if err := c.BodyParser(&body); err != nil {
		return apperr.WriteError(c, apperr.New(apperr.CodeBadRequest, "invalid request body", err))
	}
	resp, err := h.appointmentService.BookAppointment(ctx, &body)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.Created(c, resp)
}

func (h *AppointmentHandler) CreateDoctorShift(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	var body dto.CreateDoctorShiftRequest
	if err := c.BodyParser(&body); err != nil {
		return apperr.WriteError(c, apperr.New(apperr.CodeBadRequest, "invalid request body", err))
	}
	if err := h.appointmentService.CreateDoctorShift(ctx, &body); err != nil {
		return apperr.WriteError(c, err)
	}
	return response.Created(c, fiber.Map{"message": "Doctor shift created successfully"})
}

func (h *AppointmentHandler) DeleteDoctorShift(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	var body dto.DeleteDoctorShiftRequest
	if err := c.BodyParser(&body); err != nil {
		return apperr.WriteError(c, apperr.New(apperr.CodeBadRequest, "invalid request body", err))
	}
	if err := h.appointmentService.DeleteDoctorShift(ctx, body.ShiftID); err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "Doctor shift deleted successfully"})
}

func (h *AppointmentHandler) GetDoctorActiveShifts(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	shifts, err := h.appointmentService.GetDoctorActiveShifts(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, shifts)
}

func (h *AppointmentHandler) GetDoctorIncomingAppointments(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointments, err := h.appointmentService.GetDoctorIncomingAppointments(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, appointments)
}
