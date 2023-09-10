package files

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/psql"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	project := chi.URLParam(r, "project")
	fmt.Println("Project:", project)

	resp := map[string]interface{}{}
	resp["success"] = true
	resp["message"] = "projects fetched succesfully"

	db, err := psql.GetDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	schema, err := psql.SanitizeSchemaName("project_" + project)
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd := "select uuid from " + schema + ".raster_geoms;"

	rows, err := db.Query(cmd)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	uuids := []string{}
	for rows.Next() {
		var uuid string
		err = rows.Scan(&uuid)
		if err != nil {
			return
		}
		uuids = append(uuids, uuid)
	}

	resp["files"] = uuids
	retcode := http.StatusOK
	w.WriteHeader(retcode)
	b, _ := json.Marshal(resp)
	w.Write(b)
}
