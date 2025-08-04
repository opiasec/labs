CREATE TABLE lab_attempts (
  id UUID PRIMARY KEY,
  external_user_id TEXT NOT NULL,
  lab_id UUID NOT NULL REFERENCES labs(id) ON DELETE CASCADE,
  namespace UUID NOT NULL, -- Unique namespace for the lab attempt
  started_at TIMESTAMP DEFAULT NOW(),
  finished_at TIMESTAMP,
  status_id UUID, -- Will add FK constraint in later migration
  rating INTEGER CHECK (rating BETWEEN 0 AND 5),
  feedback TEXT,
  score INTEGER CHECK (score BETWEEN 0 AND 100),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Create indexes for better performance
CREATE INDEX idx_lab_attempts_lab_id ON lab_attempts(lab_id);
CREATE INDEX idx_lab_attempts_status ON lab_attempts(status_id);
CREATE INDEX idx_lab_attempts_external_user_id ON lab_attempts(external_user_id);
