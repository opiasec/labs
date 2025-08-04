-- Drop Lab Transitions Table
DROP INDEX IF EXISTS idx_lab_transitions_automatic;
DROP INDEX IF EXISTS idx_lab_transitions_name;
DROP INDEX IF EXISTS idx_lab_transitions_to_status;
DROP INDEX IF EXISTS idx_lab_transitions_from_status;
DROP TABLE IF EXISTS lab_transitions;
