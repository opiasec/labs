-- Remove indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_active;

-- Drop table
DROP TABLE IF EXISTS users;
