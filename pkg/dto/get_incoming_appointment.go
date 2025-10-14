package dto

type GetIncomingAppointmentResponse struct {
	DoctorID   string               `json:"doctor_id"`
	DoctorName string               `json:"doctor_name"`
	Specialty  string               `json:"specialty"`
	StartTime  string               `json:"start_time"`
	EndTime    string               `json:"end_time"`
}