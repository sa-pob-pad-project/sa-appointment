package repository

import (
	"appointment-service/pkg/models"
	"time"

	"gorm.io/gorm"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{
		db: db,
	}
}

func (r *AppointmentRepository) GetDoctorShifts(doctorID string) (*[]models.DoctorShift, error) {
	var shifts []models.DoctorShift
	if err := r.db.Where("doctor_id = ? AND deleted_at IS NULL", doctorID).Find(&shifts).Error; err != nil {
		return nil, err
	}
	return &shifts, nil
}

func (r *AppointmentRepository) GetAppointmentsOfDoctor(doctorID string, startDate, endDate time.Time) (*[]models.Appointment, error) {
	var appointments []models.Appointment
	if err := r.db.Where("doctor_id = ? AND start_time >= ? AND end_time <= ? AND deleted_at IS NULL", doctorID, startDate, endDate).Find(&appointments).Error; err != nil {
		return nil, err
	}
	return &appointments, nil
}

func (r *AppointmentRepository) CreateAppointment(appointment *models.Appointment) error {
	if err := r.db.Create(appointment).Error; err != nil {
		return err
	}
	return nil
}

func (r *AppointmentRepository) CreateDoctorShift(shift *models.DoctorShift) error {
	if err := r.db.Create(shift).Error; err != nil {
		return err
	}
	return nil
}
