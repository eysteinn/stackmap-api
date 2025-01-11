
--- BEGIN OF USER
CREATE TABLE IF NOT EXISTS public.users (
    user_id SERIAL PRIMARY KEY,         -- Auto-incremented unique user ID
    --username VARCHAR(50) NOT NULL UNIQUE, -- Unique username for login
    email VARCHAR(100) NOT NULL UNIQUE, -- Unique email address
    password_hash TEXT NOT NULL,        -- Hashed password for authentication
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of user creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp of the last update
);

-- Create a trigger function to update the "updated_at" column
CREATE OR REPLACE FUNCTION update_updated_at_column_users()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger function to the "users" table
DROP TRIGGER IF EXISTS set_updated_at_column_user_trigger ON public.users;
CREATE TRIGGER set_updated_at_column_user_trigger BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column_users();
--- END OF USER



--- BEGIN OF USER_IDS: Use for external logins
CREATE TABLE IF NOT EXISTS public.user_ids (
    id SERIAL PRIMARY KEY,             -- Auto-incremented unique ID
    user_id INT NOT NULL,              -- Foreign key referencing users table
    id_type VARCHAR(50) NOT NULL,      -- Type of ID (e.g., 'google', 'github')
    external_id VARCHAR(100) NOT NULL, -- External ID for the given type
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of ID creation
    UNIQUE (user_id, id_type),         -- Ensure one ID per type for each user
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
--- END OF USER_IDS

