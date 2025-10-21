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

// GetPatientAppointmentHistory godoc
// @Summary Get patient appointment history
// @Description Retrieves all past appointments for the authenticated patient
// @Tags Patient Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.GetAppointmentHistoryResponse} "List of patient's past appointments"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/patient/history [get]
func (h *AppointmentHandler) GetPatientAppointmentHistory(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointments, err := h.appointmentService.GetPatientAppointmentHistory(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, appointments)
}

// GetLatestAppointmentHistory godoc
// @Summary Get latest appointment history
// @Description Retrieves the most recent appointment for the authenticated patient
// @Tags Patient Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.GetAppointmentHistoryResponse} "Latest appointment details"
// @Success 204 "No appointment found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/history/latest [get]
func (h *AppointmentHandler) GetLatestAppointmentHistory(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointment, err := h.appointmentService.GetLatestAppointmentHistory(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	if appointment == nil {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return response.OK(c, appointment)
}

// IncomingAppointmentOfPatient godoc
// @Summary Get patient's incoming appointments
// @Description Retrieves all upcoming appointments for the authenticated patient
// @Tags Patient Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.GetIncomingAppointmentResponse} "List of patient's upcoming appointments"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/patient/incoming [get]
func (h *AppointmentHandler) IncomingAppointmentOfPatient(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointments, err := h.appointmentService.GetPatientIncomingAppointments(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, appointments)
}

// GetDoctorSlots godoc
// @Summary Get available doctor slots
// @Description Retrieves all available appointment slots for a specific doctor
// @Tags Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param doctor_id path string true "Doctor UUID"
// @Success 200 {object} response.Response{data=dto.GetDoctorSlotResponse} "Map of dates to available slots"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Doctor not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/doctor/{doctor_id}/slots [get]
func (h *AppointmentHandler) GetDoctorSlots(c *fiber.Ctx) error {
	doctorID := c.Params("doctor_id")
	slots, err := h.appointmentService.GetDoctorSlots(doctorID)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, slots)
}

// BookAppointment godoc
// @Summary Book an appointment
// @Description Creates a new appointment for the authenticated patient with a doctor
// @Tags Patient Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.BookAppointmentRequest true "Appointment booking details"
// @Success 201 {object} response.Response{data=dto.BookAppointmentResponse} "Appointment created successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid request body or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 409 {object} response.ErrorResponse "Conflict - slot already booked"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/patient [post]
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

// CreateDoctorShift godoc
// @Summary Create doctor shift
// @Description Creates a new recurring shift schedule for the authenticated doctor
// @Tags Doctor Shifts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateDoctorShiftRequest true "Doctor shift details"
// @Success 201 {object} response.Response{data=map[string]string} "Shift created successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid request body or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/doctor/shift [post]
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

// DeleteDoctorShift godoc
// @Summary Delete doctor shift
// @Description Soft deletes a doctor's shift schedule by shift ID
// @Tags Doctor Shifts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.DeleteDoctorShiftRequest true "Shift deletion details"
// @Success 200 {object} response.Response{data=map[string]string} "Shift deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid request body or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Shift not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/doctor/shift [delete]
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

// GetDoctorActiveShifts godoc
// @Summary Get doctor's active shifts
// @Description Retrieves all active shift schedules for the authenticated doctor
// @Tags Doctor Shifts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.GetDoctorActiveShiftsResponse} "List of doctor's active shifts"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/doctor/shift [get]
func (h *AppointmentHandler) GetDoctorActiveShifts(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	shifts, err := h.appointmentService.GetDoctorActiveShifts(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, shifts)
}

// GetDoctorIncomingAppointments godoc
// @Summary Get doctor's incoming appointments
// @Description Retrieves all upcoming appointments for the authenticated doctor
// @Tags Doctor Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.GetDoctorIncomingAppointmentsResponse} "List of doctor's upcoming appointments"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/doctor [get]
func (h *AppointmentHandler) GetDoctorIncomingAppointments(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	appointments, err := h.appointmentService.GetDoctorIncomingAppointments(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, appointments)
}

// CancelAppointment godoc
// @Summary Cancel an appointment
// @Description Cancels an existing appointment. Both patients and doctors can cancel their appointments.
// @Tags Appointments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CancelAppointmentRequest true "Appointment cancellation details"
// @Success 200 {object} response.Response{data=map[string]string} "Appointment cancelled successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid request body or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Appointment not found"
// @Failure 409 {object} response.ErrorResponse "Conflict - appointment already cancelled or in the past"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/appointment/v1/cancel [post]
func (h *AppointmentHandler) CancelAppointment(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	var body dto.CancelAppointmentRequest
	if err := c.BodyParser(&body); err != nil {
		return apperr.WriteError(c, apperr.New(apperr.CodeBadRequest, "invalid request body", err))
	}
	if err := h.appointmentService.CancelAppointment(ctx, body.AppointmentID); err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, fiber.Map{"message": "Appointment cancelled successfully"})
}
