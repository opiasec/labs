-- Create Lab Status Logs Table
-- This table tracks status changes for lab attempts, creating an audit trail

CREATE TABLE lab_status_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lab_attempt_id UUID NOT NULL,
  lab_transition_id UUID NOT NULL,
  changed_by UUID, -- Reference to user who made the change (nullable for system changes)
  comment TEXT, -- Additional information about the status change
  metadata JSONB, -- Additional metadata about the change (e.g., automated system info)
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(), -- Automatically updated on each change
  deleted_at TIMESTAMP DEFAULT NULL, -- Soft delete support
  
  -- Foreign key constraints
  CONSTRAINT fk_lab_status_logs_attempt
    FOREIGN KEY (lab_attempt_id) REFERENCES lab_attempts(id) ON DELETE CASCADE,
  CONSTRAINT fk_lab_status_logs_transition
    FOREIGN KEY (lab_transition_id) REFERENCES lab_transitions(id) ON DELETE RESTRICT
);

-- Create indexes for better performance
CREATE INDEX idx_lab_status_logs_attempt ON lab_status_logs(lab_attempt_id);
CREATE INDEX idx_lab_status_logs_transition ON lab_status_logs(lab_transition_id);
CREATE INDEX idx_lab_status_logs_changed_by ON lab_status_logs(changed_by);
CREATE INDEX idx_lab_status_logs_created_at ON lab_status_logs(created_at);

-- Create composite index for common queries (lab attempt timeline)
CREATE INDEX idx_lab_status_logs_attempt_created ON lab_status_logs(lab_attempt_id, created_at DESC);
