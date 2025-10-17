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

func (r *AppointmentRepository) GetLatestAppointmentsOfDoctor(doctorID string, startDate time.Time) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := r.db.Where("doctor_id = ? AND start_time >= ? AND deleted_at IS NULL", doctorID, startDate).Order("start_time desc").First(&appointment).Error; err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) GetAppointmentHistoryOfPatient(patientID string) (*[]models.Appointment, error) {
	var appointments []models.Appointment
	if err := r.db.Where("patient_id = ? AND deleted_at IS NULL", patientID).Order("start_time desc").Find(&appointments).Error; err != nil {
		return nil, err
	}
	return &appointments, nil
}

func (r *AppointmentRepository) GetIncomingAppointmentsOfPatient(patientID string, date time.Time) (*[]models.Appointment, error) {
	var appointments []models.Appointment
	if err := r.db.Where("patient_id = ? AND start_time >= ? AND status = 'scheduled' AND deleted_at IS NULL", patientID, date).Order("start_time asc").Find(&appointments).Error; err != nil {
		return nil, err
	}
	return &appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsOfDoctor(doctorID string, startDate, endDate time.Time) (*[]models.Appointment, error) {
	var appointments []models.Appointment
	if err := r.db.Where("doctor_id = ? AND start_time >= ? AND end_time <= ? AND deleted_at IS NULL", doctorID, startDate, endDate).Order("start_time desc").Find(&appointments).Error; err != nil {
		return nil, err
	}
	return &appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsOfPatient(patientID string) (*[]models.Appointment, error) {
	var appointments []models.Appointment
	if err := r.db.Where("patient_id = ? AND deleted_at IS NULL", patientID).Find(&appointments).Error; err != nil {
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

func (r *AppointmentRepository) GetDoctorShiftsByDoctorID(doctorID string) (*[]models.DoctorShift, error) {
	var shifts []models.DoctorShift
	if err := r.db.Where("doctor_id = ? AND deleted_at IS NULL", doctorID).Find(&shifts).Error; err != nil {
		return nil, err
	}
	return &shifts, nil
}

func (r *AppointmentRepository) GetDoctorShiftByID(shiftID string) (*models.DoctorShift, error) {
	var shift models.DoctorShift
	if err := r.db.Where("id = ? AND deleted_at IS NULL", shiftID).First(&shift).Error; err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *AppointmentRepository) DeleteDoctorShift(shiftID string) error {
	if err := r.db.Where("id = ?", shiftID).Delete(&models.DoctorShift{}).Error; err != nil {
		return err
	}
	return nil
}
