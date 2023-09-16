package products

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/psql"
)

func GetProducts(project string) (map[string]Product, error) {
	db, err := psql.GetDB()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	schema, err := psql.SanitizeSchemaName("project_" + project)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cmd := "select uuid from " + schema + ".raster_geoms;"

	//cmd := "select regexp_replace(n1.schema_name, '^project_', '') as project from (select schema_name from information_schema.schemata where schema_name ~ '^project*') n1;"

	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := map[string]Product{}
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.Name)
		if err != nil {
			return nil, err
		}
		//projects = append(projects, p)
		projects[p.Name] = p
	}
	return projects, nil

}

func GetHandler(w http.ResponseWriter, r *http.Request) {

	project := chi.URLParam(r, "project")
	log.Println("GetHandler project:", project)

	resp := map[string]interface{}{}
	resp["success"] = true
	resp["message"] = "products fetched succesfully"
	w.Header().Set("Content-Type", "application/json")
	retcode := http.StatusOK

	products, err := GetProducts(project)
	if err != nil {
		retcode = http.StatusInternalServerError
		resp["message"] = "Internal error"
		resp["success"] = false
	} else {
		resp["products"] = products
	}

	w.WriteHeader(retcode)
	b, _ := json.Marshal(resp)
	w.Write(b)
}
