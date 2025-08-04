CREATE TABLE lab_vulnerabilities (
  lab_id UUID REFERENCES labs(id) ON DELETE CASCADE,
  vulnerability_id UUID REFERENCES vulnerabilities(id) ON DELETE CASCADE,
  PRIMARY KEY (lab_id, vulnerability_id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);
