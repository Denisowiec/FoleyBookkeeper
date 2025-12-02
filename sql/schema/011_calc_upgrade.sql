-- +goose Up
ALTER TABLE calculations ADD boss_tribute NUMERIC NOT NULL DEFAULT 30;
ALTER TABLE calculations ADD manager_commission NUMERIC NOT NULL DEFAULT 3;
ALTER TABLE calculations ADD tax_rate NUMERIC NOT NULL DEFAULT 0.12;
ALTER TABLE calculations ADD tax_multiplier NUMERIC NOT NULL DEFAULT 0.5;

-- +goose Down
ALTER TABLE calculations DROP COLUMN boss_tribute;
ALTER TABLE calculations DROP COLUMN manager_commission;
ALTER TABLE calculations DROP COLUMN tax_rate;
ALTER TABLE calculations DROP COLUMN tax_multiplier;