CREATE TABLE "{{SCHEMA}}"."files" (
    file_id bigserial PRIMARY KEY,
    "filename" text NOT NULL,
    "metadata" jsonb,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX files_idx_filename ON "{{SCHEMA}}"."files" (filename);

/*CREATE TABLE "{{SCHEMA}}"."files" (gid serial, "uuid" uuid default public.uuid_generate_v4(), "filename" varchar(254), "metadata" json);
ALTER TABLE "{{SCHEMA}}"."files" ADD PRIMARY KEY (gid);
/* Create unique index on uuid, also this is a foreign key in raster_geoms */
CREATE UNIQUE INDEX files_idx_uuid on "{{SCHEMA}}"."files"(uuid); 
CREATE INDEX files_idx_filename on "{{SCHEMA}}"."files"(filename); */
/* insert json: insert into "project_vedur"."files" (filename, metadata) values ('test', '{ "rugl": "abc" }'); */