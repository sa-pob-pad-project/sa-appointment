package dto

import "time"

type DoctorSlot struct {
	DoctorID  string `json:"doctor_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}

type GetDoctorSlotResponse map[time.Time][]DoctorSlot
