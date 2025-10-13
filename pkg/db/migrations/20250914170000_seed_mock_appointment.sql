-- +goose Up
INSERT INTO appointments (patient_id, doctor_id, start_time, end_time, status) VALUES
  ('a3f5d6e2-1234-4b9a-8cde-0f1a2b3c4d5e', '0199474d-87ba-7a66-901a-641dcc6e8e79', '2025-09-14 10:00:00+07', '2025-09-14 11:00:00+07', 'scheduled');

-- +goose Down
DELETE FROM appointments
WHERE patient_id = 'a3f5d6e2-1234-4b9a-8cde-0f1a2b3c4d5e'
  AND doctor_id = '0199474d-87ba-7a66-901a-641dcc6e8e79'
  AND start_time = '2025-09-14 10:00:00+07';