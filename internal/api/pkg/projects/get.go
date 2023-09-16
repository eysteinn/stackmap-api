package projects

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/psql"
)

func GetProjects() (map[string]Project, error) {
	db, err := psql.GetDB()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cmd := "select regexp_replace(n1.schema_name, '^project_', '') as project from (select schema_name from information_schema.schemata where schema_name ~ '^project*') n1;"

	rows, err := db.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := map[string]Project{}
	for rows.Next() {
		var p Project
		err = rows.Scan(&p.Name)
		if err != nil {
			return nil, err
		}
		projects[p.Name] = p
		//projects = append(projects, p)
	}
	return projects, nil
}

func GetHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Handling get project request")
	resp := map[string]interface{}{}
	resp["success"] = true
	resp["message"] = "projects fetched succesfully"
	retcode := http.StatusOK
	w.Header().Set("Content-Type", "application/json")

	projects, err := GetProjects()
	if err != nil {
		retcode = http.StatusBadRequest
		resp["success"] = false
		resp["message"] = "failed to get projects"
	}
	if projects != nil {
		resp["projects"] = projects
		keys := make([]string, 0, len(projects))
		for k := range projects {
			keys = append(keys, k)
		}
		resp["project_names"] = keys
	}

	w.WriteHeader(retcode)
	b, _ := json.Marshal(resp)
	w.Write(b)
}
