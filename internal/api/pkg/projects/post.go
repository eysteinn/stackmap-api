package projects

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "github.com/lib/pq"

	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/requests"
)

var (
	db *sql.DB = nil
)

func setup() (*sql.DB, error) {

	Host := os.Getenv("PSQLHOST")
	User := os.Getenv("PSQLUSER")
	DB := os.Getenv("PSQLDB")
	Pass := os.Getenv("PSQLPASS")
	Port := os.Getenv("PSQLPORT")

	if Host == "" {
		Host = "postgresql.default.svc.cluster.local"
	}
	if User == "" {
		User = "postgres"
	}
	if DB == "" {
		DB = "postgres"
	}

	if Host == "" || User == "" || DB == "" || Pass == "" {
		return nil, errors.New("Unable to grap credentials for PSQL")
	}
	if Port == "" {
		Port = "5432"
	}

	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Host, User, Pass, DB, Port)

	return sql.Open("postgres", conn)
}

func GetDB() (*sql.DB, error) {
	if db == nil {
		tmp, err := setup()
		if err != nil {
			return nil, err
		}
		db = tmp
	}
	return db, nil
}

func sanitizeSchemaName(name string) (string, error) {
	// Replace spaces with underscores
	newname := strings.ReplaceAll(name, " ", "_")

	// Remove any characters that are not alphanumeric or underscores
	validName := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	newname = validName.ReplaceAllString(newname, "")

	// Truncate the name to a reasonable length if needed
	if len(newname) > 63 {
		newname = newname[:63]
	}

	if newname != name {
		return newname, fmt.Errorf("Project name is not allowed, use only alphanumerical characters")
	}
	return newname, nil
}
func CreateProject(project Project) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	schema, err := sanitizeSchemaName("project_" + project.Name)
	if err != nil {
		return err
	}

	create_tbl_cmd := `
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
ANALYZE "raster_geoms";`

	create_tbl_cmd = strings.ReplaceAll(create_tbl_cmd, "{SCHEMA}", schema)
	fmt.Println("Create table: ", create_tbl_cmd)
	_, err = db.Exec(create_tbl_cmd)
	if err != nil {
		return err
	}

	fmt.Println("Finished creating")
	return nil
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Post request")
	contenttype := requests.GetContent(r)
	name := ""
	switch contenttype {
	case "multipart/form-data":
		r.ParseMultipartForm(10 << 20) // 10 MB
		for k, v := range r.PostForm {
			if k == "name" && len(v) > 0 {
				name = v[0]
			}
		}
	}

	resp := map[string]interface{}{}
	resp["success"] = true
	resp["message"] = "layer created succesfully"
	w.Header().Set("Content-Type", "application/json")

	if name == "" {
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
		fmt.Println(err)
		resp["success"] = false
		resp["message"] = fmt.Sprint("error creating project")
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	resp["project"] = project
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
}
