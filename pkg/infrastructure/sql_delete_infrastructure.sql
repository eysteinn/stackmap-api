DROP TABLE IF EXISTS public.projects, public.users, public.user_projects, public.user_ids CASCADE;

DROP FUNCTION IF EXISTS update_updated_at_column_users, update_updated_at_column_projects;

-- Drop all project schemas
DO $$ 
DECLARE 
    schema_name TEXT;
BEGIN
    FOR schema_name IN 
        SELECT s.schema_name  -- Qualify with alias to avoid ambiguity
        FROM information_schema.schemata AS s
        WHERE s.schema_name ~ '^project_*'
    LOOP
        -- Drop the schema dynamically
        EXECUTE format('DROP SCHEMA IF EXISTS %I CASCADE', schema_name);
    END LOOP;
END $$;
