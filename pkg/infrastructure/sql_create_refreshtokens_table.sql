CREATE TABLE IF NOT EXISTS public.refresh_tokens (
    id SERIAL PRIMARY KEY,           -- Optional unique identifier for the token entry
    token TEXT NOT NULL UNIQUE,      -- The refresh token value, must be unique
    user_id INT NOT NULL,            -- Foreign key linking to the users table
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the token was issued
    expires_at TIMESTAMP NOT NULL,   -- When the token expires
    revoked BOOLEAN DEFAULT FALSE,   -- Whether the token has been revoked
    FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE
);
