CREATE TABLE IF NOT EXISTS {{SCHEMA}}.files_to_process (
    id SERIAL PRIMARY KEY,
    file_path TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of project creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp of the last update,
);

-- Create a trigger function to update the "updated_at" column
CREATE OR REPLACE FUNCTION {{SCHEMA}}.update_updated_at_column_files_to_process()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger function to the "files_tp_process" table
DROP TRIGGER IF EXISTS set_updated_at_files_to_process ON {{SCHEMA}}.files_to_process;
CREATE TRIGGER set_updated_at_files_to_process
BEFORE UPDATE ON {{SCHEMA}}.files_to_process
FOR EACH ROW
EXECUTE FUNCTION {{SCHEMA}}.update_updated_at_column_files_to_process();


-- Create a trigger function to notify of a new file
CREATE OR REPLACE FUNCTION {{SCHEMA}}.notify_new_file_to_process()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify(
        'file_channel',
        json_build_object(
            'schema', current_schema(),
            'id', NEW.id,
            'file_path', NEW.file_path
        )::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS file_to_process_insert_trigger ON {{SCHEMA}}.files_to_process;
CREATE TRIGGER file_to_process_insert_trigger
AFTER INSERT ON {{SCHEMA}}.files_to_process
FOR EACH ROW
EXECUTE FUNCTION {{SCHEMA}}.notify_new_file_to_process();