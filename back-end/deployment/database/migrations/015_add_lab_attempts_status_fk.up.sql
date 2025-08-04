-- Add Foreign Key Constraint to lab_attempts.status_id
-- This migration adds the FK constraint after lab_statuses table is created

ALTER TABLE lab_attempts
ADD CONSTRAINT fk_lab_attempts_status
FOREIGN KEY (status_id) REFERENCES lab_statuses(id) ON DELETE RESTRICT;
