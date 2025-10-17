-- +goose Up
-- +goose StatementBegin
INSERT INTO doctor_shifts (doctor_id, weekday, start_time, end_time, duration_min) VALUES
	('01920e5a-1234-7890-abcd-000000000006', 'mon', '2025-09-14 09:00:00+07', '2025-09-14 17:00:00+07', 60),
	('01920e5a-1234-7890-abcd-000000000005', 'tue', '2025-09-15 09:00:00+07', '2025-09-15 17:00:00+07', 60);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM doctor_shifts
WHERE doctor_id IN ('01920e5a-1234-7890-abcd-000000000006', '01920e5a-1234-7890-abcd-000000000005');
-- +goose StatementEnd
