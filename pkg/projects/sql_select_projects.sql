SELECT regexp_replace(n1.schema_name, '^project_', '') AS project FROM (SELECT schema_name FROM information_schema.schemata WHERE schema_name ~ '^project_*') n1;
