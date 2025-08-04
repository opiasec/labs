-- Remove Lab Transitions Seeds
DELETE FROM lab_transitions WHERE name IN (
  'pass_lab', 'fail_lab', 'abandon_lab', 'submit_for_evaluation',
  'auto_pass', 'auto_fail', 'require_manual_review', 'approve_lab',
  'reject_lab', 'retry_lab', 'resume_lab', 'resubmit_lab', 'restart_lab'
);
