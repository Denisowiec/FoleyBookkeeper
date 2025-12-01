-- +goose Up
ALTER TABLE calculations ADD boss_tribute NUMERIC NOT NULL DEFAULT 30;
ALTER TABLE calculations ADD manager_commission NUMERIC NOT NULL DEFAULT 3;

-- +goose Down
ALTER TABLE calculations DROP COLUMN boss_tribute;
ALTER TABLE calculations DROP COLUMN manager_commission;