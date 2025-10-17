package dto

type GetDoctorIncomingAppointmentsResponse struct {
	AppointmentID    string `json:"appointment_id"`
	PatientID        string `json:"patient_id"`
	PatientFirstName string `json:"patient_first_name"`
	PatientLastName  string `json:"patient_last_name"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Status           string `json:"status"`
}
