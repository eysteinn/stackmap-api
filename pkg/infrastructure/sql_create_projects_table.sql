
--- BEGIN OF PROJECTS
CREATE TABLE IF NOT EXISTS public.projects (
    project_id SERIAL PRIMARY KEY,                -- Auto-incremented unique project ID
    name VARCHAR(50) NOT NULL UNIQUE,               -- Name of project
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of project creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp of the last update,
    
);

-- Create a trigger function to update the "updated_at" column
CREATE OR REPLACE FUNCTION update_updated_at_column_projects()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger function to the "projects" table
DROP TRIGGER IF EXISTS set_updated_at_projects ON public.projects;
CREATE TRIGGER set_updated_at_projects
BEFORE UPDATE ON public.projects
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column_projects();
--- END OF PROJECTS

