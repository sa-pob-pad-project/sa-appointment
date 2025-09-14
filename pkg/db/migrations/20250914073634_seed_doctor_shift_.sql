-- +goose Up
-- +goose StatementBegin
INSERT INTO doctor_shifts (doctor_id, weekday, start_time, end_time, duration_min) VALUES
	('0199474d-87ba-7a66-901a-641dcc6e8e79', 'mon', '2025-09-14 09:00:00+07', '2025-09-14 17:00:00+07', 60),
	('0199474d-b649-7a9d-b13c-5de57528832f', 'tue', '2025-09-15 09:00:00+07', '2025-09-15 17:00:00+07', 60);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM doctor_shifts
WHERE doctor_id IN ('0199474d-87ba-7a66-901a-641dcc6e8e79', '0199474d-b649-7a9d-b13c-5de57528832f');
-- +goose StatementEnd
