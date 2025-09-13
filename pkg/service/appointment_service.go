package service

import (
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/repository"
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
