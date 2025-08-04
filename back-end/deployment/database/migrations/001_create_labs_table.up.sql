CREATE TABLE labs (
  id UUID PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  authors TEXT[],  
  external_references TEXT[],
  active BOOLEAN DEFAULT TRUE,
  estimated_time INT,
  requires_manual_review BOOLEAN DEFAULT FALSE, 
  description TEXT NOT NULL,
  readme TEXT, 
  difficulty TEXT CHECK (difficulty IN ('easy', 'medium', 'hard')),
  tags TEXT[],
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);
