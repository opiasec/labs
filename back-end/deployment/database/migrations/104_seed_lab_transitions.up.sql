-- Seed Lab Transitions
-- This file seeds the lab_transitions table with default transitions for the FSM

INSERT INTO lab_transitions (from_status_id, to_status_id, name, is_automatic, requires_approval)
SELECT 
  fs.id as from_status_id,
  ts.id as to_status_id,
  t.name,
  t.is_automatic,
  t.requires_approval
FROM 
  (VALUES 
    ('running', 'passed', 'pass_lab', FALSE, FALSE),
    ('running', 'failed', 'fail_lab', FALSE, FALSE),
    ('running', 'abandoned', 'abandon_lab', FALSE, FALSE),
    ('running', 'evaluating', 'submit_for_evaluation', FALSE, FALSE),
    ('evaluating', 'passed', 'auto_pass', TRUE, FALSE),
    ('evaluating', 'failed', 'auto_fail', TRUE, FALSE),
    ('evaluating', 'pending_review', 'require_manual_review', TRUE, FALSE),
    ('pending_review', 'approved', 'approve_lab', FALSE, TRUE),
    ('pending_review', 'rejected', 'reject_lab', FALSE, TRUE),
    ('failed', 'running', 'retry_lab', FALSE, FALSE),
    ('abandoned', 'running', 'resume_lab', FALSE, FALSE),
    ('rejected', 'running', 'resubmit_lab', FALSE, FALSE),
    ('approved', 'running', 'restart_lab', FALSE, FALSE)
  ) as t(from_status, to_status, name, is_automatic, requires_approval)
  JOIN lab_statuses fs ON fs.name = t.from_status
  JOIN lab_statuses ts ON ts.name = t.to_status;
