CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  name TEXT NOT NULL,
  image_url TEXT,
  role TEXT DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
  is_active BOOLEAN DEFAULT TRUE,
  email_verified BOOLEAN DEFAULT FALSE,
  last_login_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Index for better performance on email lookups
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;

-- Index for role-based queries
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;

-- Index for active users
CREATE INDEX idx_users_active ON users(is_active) WHERE deleted_at IS NULL;