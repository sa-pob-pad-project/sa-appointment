package dto

import "time"

// type DayOfWeek string

// const (
// 	DayMon DayOfWeek = "mon"
// 	DayTue DayOfWeek = "tue"
// 	DayWed DayOfWeek = "wed"
// 	DayThu DayOfWeek = "thu"
// 	DayFri DayOfWeek = "fri"
// 	DaySat DayOfWeek = "sat"
// 	DaySun DayOfWeek = "sun"
// )

// ID          uuid.UUID  `db:"id" json:"id"`
// DoctorID    uuid.UUID  `db:"doctor_id" json:"doctor_id"`
// Weekday     DayOfWeek  `db:"weekday" json:"weekday"`
// StartTime   time.Time  `db:"start_time" json:"start_time"`
// EndTime     time.Time  `db:"end_time" json:"end_time"`
// DurationMin int        `db:"duration_min" json:"duration_min"`
// CreatedAt   time.Time  `db:"created_at" json:"created_at"`
// UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
// DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
type CreateDoctorShiftRequest struct {
	Weekday     string    `json:"weekday" validate:"required,oneof=mon tue wed thu fri sat sun"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	DurationMin int       `json:"duration_min" validate:"required,min=1"`
}
