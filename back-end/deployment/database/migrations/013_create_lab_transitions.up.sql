-- Create Lab Transitions Table (Finite State Machine)
-- This table defines the allowed transitions between lab statuses

CREATE TABLE lab_transitions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  from_status_id UUID NOT NULL,
  to_status_id UUID NOT NULL,
  name VARCHAR(100) NOT NULL, -- Name of the transition
  is_automatic BOOLEAN DEFAULT FALSE, -- Whether this transition happens automatically
  requires_approval BOOLEAN DEFAULT FALSE, -- Whether this transition requires approval
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL, -- Soft delete support
  
  -- Foreign key constraints
  CONSTRAINT fk_lab_transitions_from_status
    FOREIGN KEY (from_status_id) REFERENCES lab_statuses(id) ON DELETE CASCADE,
  CONSTRAINT fk_lab_transitions_to_status
    FOREIGN KEY (to_status_id) REFERENCES lab_statuses(id) ON DELETE CASCADE,
  
  -- Ensure we don't have duplicate transitions
  CONSTRAINT unique_transition UNIQUE (from_status_id, to_status_id, name)
);

-- Create indexes for better performance
CREATE INDEX idx_lab_transitions_from_status ON lab_transitions(from_status_id);
CREATE INDEX idx_lab_transitions_to_status ON lab_transitions(to_status_id);
CREATE INDEX idx_lab_transitions_name ON lab_transitions(name);
CREATE INDEX idx_lab_transitions_automatic ON lab_transitions(is_automatic);
