package dto

type CancelAppointmentRequest struct {
	AppointmentID string `json:"appointment_id" validate:"required,uuid"`
}

type CancelAppointmentResponse struct {
	Message string `json:"message"`
}
