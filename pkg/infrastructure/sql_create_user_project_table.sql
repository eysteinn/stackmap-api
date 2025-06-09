
--- BEGIN OF USER_PROJECTS: Used to link users and projects
CREATE TABLE IF NOT EXISTS public.user_projects (
    id SERIAL PRIMARY KEY,           -- Optional primary key for the table
    user_id INT NOT NULL,            -- Foreign key to users table
    project_id INT NOT NULL,         -- Foreign key to projects table
    role VARCHAR(50) DEFAULT 'member', -- Role or type of access (e.g., 'admin', 'viewer')
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of when the association was created
    UNIQUE (user_id, project_id),    -- Prevents duplicate associations
    FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES public.projects(project_id) ON DELETE CASCADE
);