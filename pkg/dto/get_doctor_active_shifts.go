package dto

type GetDoctorActiveShiftsResponse struct {
	ShiftID     string `json:"shift_id"`
	Weekday     string `json:"weekday"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	DurationMin int    `json:"duration_min"`
}
