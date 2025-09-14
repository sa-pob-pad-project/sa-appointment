package service

import (
	"appointment-service/pkg/dto"
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/models"
	"appointment-service/pkg/repository"
	"fmt"
	"strings"
	"time"
)

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	jwtService      *jwt.JwtService
}

func NewAppointmentService(appointmentRepo *repository.AppointmentRepository, jwtService *jwt.JwtService) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		jwtService:      jwtService,
	}
}

func (s *AppointmentService) GetDoctorSlots(doctorID string) (*dto.GetDoctorSlotResponse, error) {

	shifts, err := s.appointmentRepo.GetDoctorShifts(doctorID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	appointments, err := s.appointmentRepo.GetAppointmentsOfDoctor(doctorID, startDate, startDate.AddDate(0, 0, 30))
	if err != nil {
		return nil, err
	}
	slots := map[time.Time][]dto.DoctorSlot{}
	// array of day from now to now + 30 days
	for i := 0; i < 30; i++ {
		day := startDate.AddDate(0, 0, i)
		weekDay := strings.ToLower(day.Weekday().String()[:3])
		var shift *models.DoctorShift
		for _, s := range *shifts {
			if string(s.Weekday) == weekDay {
				shift = &s
			}
		}

		if shift == nil {
			continue
		}

		slotInterval := time.Duration(shift.DurationMin) * time.Minute // convert minutes to Duration
		slots[day] = []dto.DoctorSlot{}

		for t := shift.StartTime; t.Before(shift.EndTime); t = t.Add(slotInterval) {
			slots[day] = append(slots[day], dto.DoctorSlot{
				DoctorID:  doctorID,
				StartTime: t.Format("15:04Z00:00"),                   // include timezone offset
				EndTime:   t.Add(slotInterval).Format("15:04Z00:00"), // include timezone offset
				Status:    "available",
			})
		}
	}

	fmt.Println("Appointments", appointments)

	return (*dto.GetDoctorSlotResponse)(&slots), nil
}
