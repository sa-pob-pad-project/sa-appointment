package dto

type DeleteDoctorShiftRequest struct {
	ShiftID string `json:"shift_id" validate:"required,uuid"`
}

type DeleteDoctorShiftResponse struct {
	Message string `json:"message"`
}
