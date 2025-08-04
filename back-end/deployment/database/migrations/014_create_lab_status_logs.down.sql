-- Drop Lab Status Logs Table
DROP INDEX IF EXISTS idx_lab_status_logs_attempt_created;
DROP INDEX IF EXISTS idx_lab_status_logs_created_at;
DROP INDEX IF EXISTS idx_lab_status_logs_changed_by;
DROP INDEX IF EXISTS idx_lab_status_logs_to_status;
DROP INDEX IF EXISTS idx_lab_status_logs_from_status;
DROP INDEX IF EXISTS idx_lab_status_logs_transition;
DROP INDEX IF EXISTS idx_lab_status_logs_attempt;
DROP TABLE IF EXISTS lab_status_logs;
