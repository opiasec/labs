-- Create Lab Statuses Table
-- This table stores the possible statuses that a lab can have

CREATE TABLE lab_statuses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(50) UNIQUE NOT NULL,
  description TEXT,
  is_initial BOOLEAN DEFAULT FALSE, -- Indicates if this is an initial state
  is_final BOOLEAN DEFAULT FALSE, -- Indicates if this is a final state
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Create index for faster lookups
CREATE INDEX idx_lab_statuses_name ON lab_statuses(name);
CREATE INDEX idx_lab_statuses_initial ON lab_statuses(is_initial);
CREATE INDEX idx_lab_statuses_final ON lab_statuses(is_final);
