-- Remove Lab Statuses Seeds
DELETE FROM lab_statuses WHERE name IN (
  'running', 'failed', 'passed', 'evaluating', 
  'abandoned', 'pending_review', 'rejected', 'approved'
);
