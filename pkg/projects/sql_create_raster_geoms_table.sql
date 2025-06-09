SET CLIENT_ENCODING TO UTF8;
SET STANDARD_CONFORMING_STRINGS TO ON;
/*CREATE SCHEMA {{SCHEMA}};
SET SCHEMA '{{SCHEMA}}';*/
SET search_path = {{SCHEMA}}, public;
BEGIN;
CREATE TABLE {{SCHEMA}}."raster_geoms" (
    gid serial PRIMARY KEY,
    file_id INT NOT NULL,
    "location" varchar(254),
    "src_srs" varchar(254),
    "datetime" timestamp without time zone,
    "product" varchar(254)
);
/*CREATE TABLE "raster_geoms" (gid serial, file_id INT NOT NULL,            -- Foreign key to users table
 "location" varchar(254),"src_srs" varchar(254),"datetime" timestamp without time zone,"product" varchar(254));*/

/*ALTER TABLE {{SCHEMA}}."raster_geoms" ADD PRIMARY KEY (gid);
ALTER TABLE {{SCHEMA}}.raster_geoms add foreign key (file_id) references {{SCHEMA}}.files(gid) ON DELETE CASCADE;
*/
ALTER TABLE {{SCHEMA}}.raster_geoms 
ADD CONSTRAINT fk_file_id FOREIGN KEY (file_id) 
REFERENCES {{SCHEMA}}.files(file_id) ON DELETE CASCADE;


SELECT AddGeometryColumn('raster_geoms','geom','4326','MULTIPOLYGON',2);
--CREATE index idx_uuid on raster_geoms(uuid);
CREATE index idx_product_time on raster_geoms(product, datetime);

COMMIT;
ANALYZE {{SCHEMA}}."raster_geoms";