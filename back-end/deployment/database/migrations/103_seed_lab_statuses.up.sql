-- Seed Lab Statuses
-- This file seeds the lab_statuses table with default status values

INSERT INTO lab_statuses (name, description, is_initial, is_final) VALUES 
('running', 'Lab is currently being worked on', FALSE, FALSE),
('timeout', 'Lab attempt has succeeded', FALSE, TRUE),
('failed', 'Lab attempt has failed', FALSE, TRUE),
('passed', 'Lab attempt has passed', FALSE, TRUE),
('evaluating', 'Lab attempt is being evaluated', FALSE, FALSE),
('abandoned', 'Lab was abandoned by the user', FALSE, TRUE),
('pending_review', 'Lab completion is pending review', FALSE, FALSE),
('rejected', 'Lab completion has been rejected', FALSE, TRUE),
('approved', 'Lab has been reviewed and approved', FALSE, TRUE);
