package projects

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "embed"

	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/psql"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/requests"
)

//go:embed sql_create_files_table.sql
var sql_create_files_table string

//go:embed sql_create_raster_geoms_table.sql
var sql_create_raster_geoms_table string

//go:embed sql_create_schema.sql
var sql_create_schema string

func CreateProject(project Project) error {
	log.Println("Creating project: " + project.Name)
	db, err := psql.GetDB()
	if err != nil {
		return err
	}

	schema, err := psql.SanitizeSchemaName("project_" + project.Name)
	if err != nil {
		return err
	}

	/*create_tbl_cmd := `
		SET CLIENT_ENCODING TO UTF8;
	SET STANDARD_CONFORMING_STRINGS TO ON;
	CREATE SCHEMA {SCHEMA};
	SET SCHEMA '{SCHEMA}';
	SET search_path = {SCHEMA}, public;
	BEGIN;
	CREATE TABLE "raster_geoms" (gid serial, "uuid" uuid default uuid_generate_v4(), "location" varchar(254),"src_srs" varchar(254),"datetime" timestamp without time zone,"product" varchar(254));
	ALTER TABLE "raster_geoms" ADD PRIMARY KEY (gid);
	SELECT AddGeometryColumn('raster_geoms','geom','4326','MULTIPOLYGON',2);
	CREATE index idx_uuid on raster_geoms(uuid);
	CREATE index idx_product_time on raster_geoms(product, datetime);
	COMMIT;
	ANALYZE "raster_geoms";`*/

	//create_tbl_cmd = strings.ReplaceAll(create_tbl_cmd, "{SCHEMA}", schema)
	create_schema_cmd := strings.ReplaceAll(sql_create_schema, "{{SCHEMA}}", schema)
	fmt.Println("Create schema: ", create_schema_cmd)
	_, err = db.Exec(create_schema_cmd)
	if err != nil {
		return err
	}

	create_files_cmd := strings.ReplaceAll(sql_create_files_table, "{{SCHEMA}}", schema)
	_, err = db.Exec(create_files_cmd)
	fmt.Println("Finished creating")

	create_raster_geoms_cmd := strings.ReplaceAll(sql_create_raster_geoms_table, "{{SCHEMA}}", schema)
	fmt.Println("Create table: ", create_raster_geoms_cmd)
	_, err = db.Exec(create_raster_geoms_cmd)
	if err != nil {
		return err
	}
	return nil

}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Post request")
	contenttype := requests.GetContent(r)
	name := ""
	resp := map[string]interface{}{}
	resp["success"] = true
	resp["message"] = "layer created succesfully"
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Incoming request:", contenttype)
	switch contenttype {
	case "multipart/form-data":
		r.ParseMultipartForm(10 << 20) // 10 MB
		for k, v := range r.PostForm {
			fmt.Println(k, "=", v)
			if k == "name" && len(v) > 0 {
				name = v[0]
			}
		}
	case "application/x-www-form-urlencoded":
		r.ParseForm()
		for k, v := range r.PostForm {
			fmt.Println(k, "=", v)
			if k == "name" && len(v) > 0 {
				name = v[0]
			}
		}
	default:

		resp["success"] = false
		resp["message"] = "ContentType should be multipart/form-data, got " + contenttype + " instead."
		fmt.Println(resp["message"])
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	if name == "" {
		log.Println("Missing name")
		resp["success"] = false
		resp["message"] = "missing 'name' parameter"
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	project := Project{
		Name: name,
	}

	err := CreateProject(project)
	if err != nil {
		log.Println(err)
		resp["success"] = false
		resp["message"] = fmt.Sprint("error creating project")
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	log.Println("Project created successfully")
	resp["project"] = project
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
}
