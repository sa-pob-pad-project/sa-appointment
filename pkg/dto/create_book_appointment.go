package dto

import "time"

type BookAppointmentRequest struct {
	DoctorID  string    `json:"doctor_id" validate:"required,uuid"`
	StartTime time.Time `json:"start_time" validate:"required"`
}

type BookAppointmentResponse struct {
	AppointmentID string `json:"appointment_id"`
}
