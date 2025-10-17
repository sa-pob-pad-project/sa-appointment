package service

import (
	"appointment-service/pkg/apperr"
	"appointment-service/pkg/clients"
	client_dto "appointment-service/pkg/clients/dto"
	contextUtils "appointment-service/pkg/context"
	"appointment-service/pkg/dto"
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/models"
	"appointment-service/pkg/repository"
	"appointment-service/pkg/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	userClient      *clients.UserClient
	jwtService      *jwt.JwtService
}

func NewAppointmentService(appointmentRepo *repository.AppointmentRepository, userClient *clients.UserClient, jwtService *jwt.JwtService) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		userClient:      userClient,
		jwtService:      jwtService,
	}
}

func (s *AppointmentService) GetPatientAppointmentHistory(ctx context.Context) (*[]dto.GetAppointmentHistoryResponse, error) {
	patientID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)
	if role != "patient" {
		return nil, apperr.New(apperr.CodeForbidden, "only patients can view their appointment history", nil)
	}
	appointments, err := s.appointmentRepo.GetAppointmentHistoryOfPatient(patientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &[]dto.GetAppointmentHistoryResponse{}, nil
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch appointments", err)
	}
	doctorIds := []string{}
	for _, appointment := range *appointments {
		doctorIds = append(doctorIds, appointment.DoctorID.String())
	}
	fmt.Println("appointment", appointments)
	if len(doctorIds) == 0 {
		return &[]dto.GetAppointmentHistoryResponse{}, nil
	}
	doctorProfiles, err := s.userClient.GetDoctorByIds(ctx, doctorIds)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor profile "+err.Error(), err)
	}
	doctorIdToProfile := map[string]client_dto.GetDoctorProfileResponseDto{}
	for _, profile := range *doctorProfiles {
		doctorIdToProfile[profile.ID] = profile
	}
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor profile", err)
	}
	var responses []dto.GetAppointmentHistoryResponse
	for _, appointment := range *appointments {
		doctorProfile := doctorIdToProfile[appointment.DoctorID.String()]
		responses = append(responses, dto.GetAppointmentHistoryResponse{
			DoctorID:        appointment.DoctorID.String(),
			DoctorFirstName: doctorProfile.FirstName,
			DoctorLastName:  doctorProfile.LastName,
			Specialty:       *doctorProfile.Specialty,
			StartTime:       appointment.StartTime.Format("2006-01-02 15:04Z07:00"),
			EndTime:         appointment.EndTime.Format("2006-01-02 15:04Z07:00"),
			Status:          string(appointment.Status),
		})
	}
	return &responses, nil
}

func (s *AppointmentService) GetPatientIncomingAppointments(ctx context.Context) (*[]dto.GetIncomingAppointmentResponse, error) {
	patientID := contextUtils.GetUserId(ctx)
	appointments, err := s.appointmentRepo.GetIncomingAppointmentsOfPatient(patientID, time.Now())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &[]dto.GetIncomingAppointmentResponse{}, nil
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch appointments", err)
	}
	doctorIds := []string{}
	for _, appointment := range *appointments {
		doctorIds = append(doctorIds, appointment.DoctorID.String())
	}
	doctorProfiles, err := s.userClient.GetDoctorByIds(ctx, doctorIds)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor profile", err)
	}

	var responses []dto.GetIncomingAppointmentResponse
	for i, appointment := range *appointments {
		doctorProfile := (*doctorProfiles)[i]
		responses = append(responses, dto.GetIncomingAppointmentResponse{
			DoctorID:        appointment.DoctorID.String(),
			DoctorFirstName: doctorProfile.FirstName,
			DoctorLastName:  doctorProfile.LastName,
			Specialty:       *doctorProfile.Specialty,
			StartTime:       appointment.StartTime.Format("2006-01-02 15:04Z07:00"),
			EndTime:         appointment.EndTime.Format("2006-01-02 15:04Z07:00"),
		})
	}
	return &responses, nil
}

func (s *AppointmentService) GetDoctorSlots(doctorID string) (*dto.GetDoctorSlotResponse, error) {

	shifts, err := s.appointmentRepo.GetDoctorShifts(doctorID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor shifts", err)
	}
	for _, shift := range *shifts {
		fmt.Println("Shift", shift)
	}
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	appointments, err := s.appointmentRepo.GetAppointmentsOfDoctor(doctorID, startDate, startDate.AddDate(0, 0, 30))
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor appointments", err)
	}
	appointmentsMap := map[string]models.Appointment{}
	for _, app := range *appointments {
		appointmentsMap[app.StartTime.Format("2006-01-02 15:04Z00:00")] = app
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

func (s *AppointmentService) BookAppointment(ctx context.Context, req *dto.BookAppointmentRequest) (*dto.BookAppointmentResponse, error) {
	patientID := contextUtils.GetUserId(ctx)
	startTime := req.StartTime

	endTime := startTime.Add(60 * time.Minute)
	appointmentID := utils.GenerateUUIDv7()
	doctorID := utils.StringToUUIDv7(req.DoctorID)
	patientUUID := utils.StringToUUIDv7(patientID)
	_, err := s.userClient.GetDoctorByIds(ctx, []string{req.DoctorID})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperr.New(apperr.CodeNotFound, "doctor not found", err)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor profile", err)
	}
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
		return nil, apperr.New(apperr.CodeInternal, "failed to create appointment", err)
	}
	return &dto.BookAppointmentResponse{
		AppointmentID: appointmentID.String(),
	}, nil

}

func (s *AppointmentService) CreateDoctorShift(ctx context.Context, shift *dto.CreateDoctorShiftRequest) error {
	role := contextUtils.GetRole(ctx)
	if role != "doctor" {
		return apperr.New(apperr.CodeForbidden, "only doctors can create shifts", nil)
	}
	doctorID := contextUtils.GetUserId(ctx)
	shiftModel := &models.DoctorShift{
		ID:          utils.GenerateUUIDv7(),
		DoctorID:    utils.StringToUUIDv7(doctorID),
		Weekday:     models.DayOfWeek(strings.ToLower(shift.Weekday)),
		StartTime:   shift.StartTime,
		EndTime:     shift.EndTime,
		DurationMin: shift.DurationMin,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.appointmentRepo.CreateDoctorShift(shiftModel); err != nil {
		return apperr.New(apperr.CodeInternal, "failed to create doctor shift", err)
	}
	return nil
}

func (s *AppointmentService) DeleteDoctorShift(ctx context.Context, shiftID string) error {
	role := contextUtils.GetRole(ctx)
	if role != "doctor" {
		return apperr.New(apperr.CodeForbidden, "only doctors can delete shifts", nil)
	}
	doctorID := contextUtils.GetUserId(ctx)
	shift, err := s.appointmentRepo.GetDoctorShiftByID(shiftID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperr.New(apperr.CodeNotFound, "shift not found", err)
		}
		return apperr.New(apperr.CodeInternal, "failed to fetch shift", err)
	}
	if shift.DoctorID.String() != doctorID {
		return apperr.New(apperr.CodeForbidden, "unauthorized to delete this shift", nil)
	}
	if err := s.appointmentRepo.DeleteDoctorShift(shiftID); err != nil {
		return apperr.New(apperr.CodeInternal, "failed to delete shift", err)
	}
	return nil
}

func (s *AppointmentService) GetDoctorActiveShifts(ctx context.Context) (*[]dto.GetDoctorActiveShiftsResponse, error) {
	role := contextUtils.GetRole(ctx)
	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can view their shifts", nil)
	}
	doctorID := contextUtils.GetUserId(ctx)
	shifts, err := s.appointmentRepo.GetDoctorShifts(doctorID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &[]dto.GetDoctorActiveShiftsResponse{}, nil
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch doctor shifts", err)
	}

	responses := []dto.GetDoctorActiveShiftsResponse{}
	for _, shift := range *shifts {
		responses = append(responses, dto.GetDoctorActiveShiftsResponse{
			ShiftID:     shift.ID.String(),
			Weekday:     string(shift.Weekday),
			StartTime:   shift.StartTime.Format("15:04Z00:00"),
			EndTime:     shift.EndTime.Format("15:04Z00:00"),
			DurationMin: shift.DurationMin,
		})
	}
	fmt.Println("Shifts", shifts)
	return &responses, nil
}

func (s *AppointmentService) GetDoctorIncomingAppointments(ctx context.Context) (*[]dto.GetDoctorIncomingAppointmentsResponse, error) {
	role := contextUtils.GetRole(ctx)
	if role != "doctor" {
		return nil, apperr.New(apperr.CodeForbidden, "only doctors can view their appointments", nil)
	}
	doctorID := contextUtils.GetUserId(ctx)
	appointments, err := s.appointmentRepo.GetIncomingAppointmentsOfDoctor(doctorID, time.Now())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &[]dto.GetDoctorIncomingAppointmentsResponse{}, nil
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch appointments", err)
	}

	patientIds := []string{}
	for _, appointment := range *appointments {
		patientIds = append(patientIds, appointment.PatientID.String())
	}

	patientProfiles, err := s.userClient.GetPatientByIds(ctx, patientIds)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to fetch patient profiles", err)
	}

	patientIdToProfile := map[string]client_dto.GetPatientProfileResponseDto{}
	for _, profile := range *patientProfiles {
		patientIdToProfile[profile.ID] = profile
	}

	var responses []dto.GetDoctorIncomingAppointmentsResponse
	for _, appointment := range *appointments {
		patientProfile := patientIdToProfile[appointment.PatientID.String()]
		responses = append(responses, dto.GetDoctorIncomingAppointmentsResponse{
			AppointmentID:    appointment.ID.String(),
			PatientID:        appointment.PatientID.String(),
			PatientFirstName: patientProfile.FirstName,
			PatientLastName:  patientProfile.LastName,
			StartTime:        appointment.StartTime.Format("2006-01-02 15:04Z07:00"),
			EndTime:          appointment.EndTime.Format("2006-01-02 15:04Z07:00"),
			Status:           string(appointment.Status),
		})
	}
	return &responses, nil
}
