package main

import (
	"appointment-service/cmd"
	"appointment-service/pkg/config"
	"appointment-service/pkg/db"
	"appointment-service/pkg/handlers"
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/repository"
	"appointment-service/pkg/routes"
	"appointment-service/pkg/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @title User API
// @description This is a sample server for a user API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	os.Setenv("TZ", "Asia/Bangkok")
	config.LoadConfig()
	db := db.Open(db.Config{
		Host:     config.Get("DB_HOST", "localhost"),
		Port:     config.GetInt("DB_PORT", 5432),
		User:     config.Get("DB_USER", "user"),
		Password: config.Get("DB_PASSWORD", "password"),
		Dbname:   config.Get("DB_NAME", "postgres"),
		Sslmode:  config.Get("DB_SSLMODE", "disable"),
	})
	cmd.InitCmd()
	appointmentRepository := repository.NewAppointmentRepository(db)
	jwtService := jwt.NewJwtService(
		config.Get("JWT_SECRET", "secret"),
		config.GetInt("JWT_TTL", 3600),
	)
	appointmentService := service.NewAppointmentService(appointmentRepository, jwtService)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	app := fiber.New(fiber.Config{
		JSONDecoder: func(b []byte, v any) error {
			dec := json.NewDecoder(bytes.NewReader(b))
			dec.DisallowUnknownFields()
			if err := dec.Decode(v); err != nil {
				return fmt.Errorf("decode: %w", err)
			}
			if err := dec.Decode(new(struct{})); err != io.EOF {
				return fmt.Errorf("decode: trailing data")
			}

			rv := reflect.ValueOf(v)
			for rv.Kind() == reflect.Pointer {
				rv = rv.Elem()
			}
			validate := validator.New()
			if rv.Kind() == reflect.Struct {
				if err := validate.Struct(v); err != nil {
					return err
				}
			}
			return nil
		},
	})

	routes.SetupRoutes(app, appointmentHandler, jwtService)
	port := config.Get("APP_PORT", "8000")
	fmt.Println("Server is running on port " + port)
	app.Listen("localhost:" + port)
}
