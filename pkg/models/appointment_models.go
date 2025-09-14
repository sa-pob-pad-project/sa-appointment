package models

import (
	"time"

	"github.com/google/uuid"
)

// SlotStatus represents the status of a doctor slot
type SlotStatus string

const (
	SlotStatusOpen      SlotStatus = "open"
	SlotStatusBooked    SlotStatus = "booked"
	SlotStatusCancelled SlotStatus = "cancelled"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

// DayOfWeek represents days of the week
type DayOfWeek string

const (
	DayMon DayOfWeek = "mon"
	DayTue DayOfWeek = "tue"
	DayWed DayOfWeek = "wed"
	DayThu DayOfWeek = "thu"
	DayFri DayOfWeek = "fri"
	DaySat DayOfWeek = "sat"
	DaySun DayOfWeek = "sun"
)

// DoctorShift maps to the doctor_shifts table
type DoctorShift struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	DoctorID    uuid.UUID  `db:"doctor_id" json:"doctor_id"`
	Weekday     DayOfWeek  `db:"weekday" json:"weekday"`
	StartTime   time.Time  `db:"start_time" json:"start_time"`
	EndTime     time.Time  `db:"end_time" json:"end_time"`
	DurationMin int        `db:"duration_min" json:"duration_min"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Appointment maps to the appointments table
type Appointment struct {
	ID        uuid.UUID         `db:"id" json:"id"`
	PatientID uuid.UUID         `db:"patient_id" json:"patient_id"`
	DoctorID  uuid.UUID         `db:"doctor_id" json:"doctor_id"`
	StartTime time.Time         `db:"start_time" json:"start_time"`
	EndTime   time.Time         `db:"end_time" json:"end_time"`
	Status    AppointmentStatus `db:"status" json:"status"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time        `db:"deleted_at" json:"deleted_at,omitempty"`
}
