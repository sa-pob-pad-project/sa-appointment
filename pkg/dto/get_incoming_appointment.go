package dto

type GetIncomingAppointmentResponse struct {
	AppointmentID   string `json:"appointment_id"`
	DoctorID        string `json:"doctor_id"`
	DoctorFirstName string `json:"doctor_first_name"`
	DoctorLastName  string `json:"doctor_last_name"`
	Specialty       string `json:"specialty"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}
