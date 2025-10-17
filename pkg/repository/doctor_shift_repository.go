package repository

import (
	"appointment-service/pkg/models"

	"gorm.io/gorm"
)

type DoctorShiftRepository struct {
	db *gorm.DB
}

func NewDoctorShiftRepository(db *gorm.DB) *DoctorShiftRepository {
	return &DoctorShiftRepository{
		db: db,
	}
}

func (r *DoctorShiftRepository) GetDoctorShifts(doctorID string) (*[]models.DoctorShift, error) {
	var shifts []models.DoctorShift
	if err := r.db.Where("doctor_id = ? AND deleted_at IS NULL", doctorID).Find(&shifts).Error; err != nil {
		return nil, err
	}
	return &shifts, nil
}

func (r *DoctorShiftRepository) CreateDoctorShift(shift *models.DoctorShift) error {
	if err := r.db.Create(shift).Error; err != nil {
		return err
	}
	return nil
}

func (r *DoctorShiftRepository) GetDoctorShiftsByDoctorID(doctorID string) (*[]models.DoctorShift, error) {
	var shifts []models.DoctorShift
	if err := r.db.Where("doctor_id = ? AND deleted_at IS NULL", doctorID).Find(&shifts).Error; err != nil {
		return nil, err
	}
	return &shifts, nil
}

func (r *DoctorShiftRepository) GetDoctorShiftByID(shiftID string) (*models.DoctorShift, error) {
	var shift models.DoctorShift
	if err := r.db.Where("id = ? AND deleted_at IS NULL", shiftID).First(&shift).Error; err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *DoctorShiftRepository) DeleteDoctorShift(shiftID string) error {
	if err := r.db.Where("id = ?", shiftID).Delete(&models.DoctorShift{}).Error; err != nil {
		return err
	}
	return nil
}
