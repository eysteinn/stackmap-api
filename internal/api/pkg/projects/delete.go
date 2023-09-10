package projects

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/psql"
)

func DeleteProject(project string) error {

	schema, err := psql.SanitizeSchemaName("project_" + project)
	if err != nil {
		return err
	}

	db, err := psql.GetDB()
	if err != nil {
		return fmt.Errorf("internal error")
	}

	_, err = db.Exec("DROP SCHEMA " + schema + " CASCADE;")
	if err != nil {
		return fmt.Errorf("unable to delete project")
	}
	return nil
}
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleting")

	project := chi.URLParam(r, "project")
	resp := map[string]interface{}{}
	w.Header().Set("Content-Type", "application/json")

	resp["success"] = true
	resp["message"] = "project deleted succesfully"
	retcode := http.StatusOK
	err := DeleteProject(project)
	if err != nil {
		resp["success"] = false
		resp["message"] = fmt.Sprint(err)
		retcode = http.StatusBadRequest
	}

	w.WriteHeader(retcode)
	b, _ := json.Marshal(resp)
	w.Write(b)

	return
}
