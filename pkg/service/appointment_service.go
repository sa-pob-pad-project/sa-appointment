package service

import (
	"appointment-service/pkg/dto"
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/models"
	"appointment-service/pkg/repository"
	"appointment-service/pkg/utils"
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
	for _, shift := range *shifts {
		fmt.Println("Shift", shift)
	}
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	appointments, err := s.appointmentRepo.GetAppointmentsOfDoctor(doctorID, startDate, startDate.AddDate(0, 0, 30))
	appointmentsMap := map[string]models.Appointment{}
	for _, app := range *appointments {
		appointmentsMap[app.StartTime.Format("2006-01-02 15:04Z00:00")] = app
	}

	if err != nil {
		return nil, err
	}
	slots := map[time.Time][]dto.DoctorSlot{}
	// array of day from now to now + 30 days
	for i := range 30 {
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
			available := "available"
			if _, exists := appointmentsMap[day.Format("2006-01-02")+" "+t.Format("15:04Z00:00")]; exists {
				available = "scheduled"

			}
			slots[day] = append(slots[day], dto.DoctorSlot{
				DoctorID:  doctorID,
				StartTime: t.Format("15:04Z00:00"),                   // include timezone offset
				EndTime:   t.Add(slotInterval).Format("15:04Z00:00"), // include timezone offset
				Status:    available,
			})
		}

	}

	fmt.Println("Appointments", appointments)

	return (*dto.GetDoctorSlotResponse)(&slots), nil
}

func (s *AppointmentService) BookAppointment(patientID string, req *dto.BookAppointmentRequest) (*dto.BookAppointmentResponse, error) {
	startTime := req.StartTime

	endTime := startTime.Add(60 * time.Minute)
	appointmentID := utils.GenerateUUIDv7()
	doctorID := utils.StringToUUIDv7(req.DoctorID)
	patientUUID := utils.StringToUUIDv7(patientID)

	appointment := &models.Appointment{
		ID:        appointmentID,
		PatientID: patientUUID,
		DoctorID:  doctorID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    "scheduled",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.appointmentRepo.CreateAppointment(appointment); err != nil {
		return nil, err
	}
	return &dto.BookAppointmentResponse{
		AppointmentID: appointmentID.String(),
	}, nil

}
