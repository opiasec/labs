-- Remove Foreign Key Constraint from lab_attempts.status_id
ALTER TABLE lab_attempts
DROP CONSTRAINT IF EXISTS fk_lab_attempts_status;
