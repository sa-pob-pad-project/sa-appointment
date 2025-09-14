-- +goose Up
-- +goose StatementBegin

CREATE TYPE slot_status AS ENUM ('open','booked','cancelled');
CREATE TYPE appointment_status AS ENUM ('scheduled','completed','cancelled');
CREATE TYPE day_of_week AS ENUM ('mon','tue','wed','thu','fri','sat','sun');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE doctor_shifts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  doctor_id uuid NOT NULL,                
  weekday day_of_week NOT NULL,
  start_time timestamptz NOT NULL,
  end_time timestamptz NOT NULL,
  duration_min int NOT NULL CHECK (duration_min > 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT shift_time_window CHECK (end_time > start_time),
  CONSTRAINT shift_unique_per_day UNIQUE (doctor_id, weekday, start_time, end_time)
);


CREATE TABLE appointments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id uuid NOT NULL,
  doctor_id uuid NOT NULL,
  start_time timestamp NOT NULL,
  end_time timestamp NOT NULL,
  status appointment_status NOT NULL DEFAULT 'scheduled',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE INDEX idx_appt_doctor_time ON appointments (doctor_id, created_at) WHERE deleted_at IS NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_appt_doctor_time;
DROP INDEX IF EXISTS idx_slots_doctor_time;

ALTER TABLE IF EXISTS appointments DROP CONSTRAINT IF EXISTS fk_appointments_slot;

-- drop tables
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS doctor_slots;
DROP TABLE IF EXISTS doctor_shifts;

-- drop extension and types
DROP EXTENSION IF EXISTS pgcrypto;
DROP TYPE IF EXISTS day_of_week;
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS slot_status;

-- drop schema
-- +goose StatementEnd