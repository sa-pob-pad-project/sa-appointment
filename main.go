package main

import (
	"appointment-service/pkg/clients"
	"appointment-service/pkg/config"
	"appointment-service/pkg/db"
	"appointment-service/pkg/handlers"
	"appointment-service/pkg/jwt"
	"appointment-service/pkg/repository"
	"appointment-service/pkg/routes"
	service "appointment-service/pkg/services"
	"bytes"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pressly/goose/v3"
)

//go:embed pkg/db/migrations/*.sql
var migrationsFS embed.FS

// migrateUp applies all up migrations using goose with embedded FS.
// Set env MIGRATE_ON_START=true (default) to run on startup.
func migrateUp(sqlDB *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}
	// directory path is relative to the root of the embedded FS
	if err := goose.Up(sqlDB, "pkg/db/migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

// @title Appointment Service API
// @description API for managing appointments, doctor shifts, and patient bookings
// @version 1.0
// @host localhost:8081
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	os.Setenv("TZ", "Asia/Bangkok")
	config.LoadConfig()
	gormDB := db.Open(db.Config{
		Host:     config.Get("DB_HOST", "localhost"),
		Port:     config.GetInt("DB_PORT", 5432),
		User:     config.Get("DB_USER", "user"),
		Password: config.Get("DB_PASSWORD", "password"),
		Dbname:   config.Get("DB_NAME", "postgres"),
		Sslmode:  config.Get("DB_SSLMODE", "disable"),
	})
	// cmd.InitCmd()
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("cannot get *sql.DB from gorm: %v", err)

	}
	if config.Get("MIGRATE_ON_START", "true") == "true" {
		if err := migrateUp(sqlDB); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}

	appointmentRepository := repository.NewAppointmentRepository(gormDB)
	doctorShiftRepository := repository.NewDoctorShiftRepository(gormDB)
	jwtService := jwt.NewJwtService(
		config.Get("JWT_SECRET", "secret"),
		config.GetInt("JWT_TTL", 3600),
	)
	userServiceUrl := config.Get("USER_SERVICE_URL", "http://localhost:8000")
	fmt.Println("userService", userServiceUrl)
	userClient := clients.New(userServiceUrl)
	appointmentService := service.NewAppointmentService(appointmentRepository, doctorShiftRepository, userClient, jwtService)
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
	// Enable CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	routes.SetupRoutes(app, appointmentHandler, jwtService)
	port := config.Get("APP_PORT", "8001")
	fmt.Println("Server is running on port " + port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
