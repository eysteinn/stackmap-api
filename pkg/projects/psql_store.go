package projects

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/lib/pq"
	"gitlab.com/EysteinnSig/stackmap-api/pkg/psql"
)

//go:embed sql_create_files_table.sql
var sql_create_files_table string

//go:embed sql_create_raster_geoms_table.sql
var sql_create_raster_geoms_table string

//go:embed sql_create_schema.sql
var sql_create_schema string

//go:embed sql_select_projects.sql
var sql_select_projects string

//go:embed sql_drop_schema.sql
var sql_drop_schema string

// PSQLProjectStore implements ProjectStore interface using an SQL database
type PSQLProjectStore struct {
	db *sql.DB
}

// NewSQLProjectStore creates a new SQLProjectStore
func NewSQLProjectStore(db *sql.DB) *PSQLProjectStore {
	return &PSQLProjectStore{db: db}
}

func (store *PSQLProjectStore) CreateProject(name, description string, ownerid int) error {
	db := store.db

	schema, err := psql.SanitizeSchemaName("project_" + name)
	if err != nil {
		return err
	}

	//create_tbl_cmd = strings.ReplaceAll(create_tbl_cmd, "{SCHEMA}", schema)
	create_schema_cmd := strings.ReplaceAll(sql_create_schema, "{{SCHEMA}}", schema)
	fmt.Println("Create schema: ", create_schema_cmd)
	_, err = db.Exec(create_schema_cmd)
	if err != nil {
		return err
	}

	create_files_cmd := strings.ReplaceAll(sql_create_files_table, "{{SCHEMA}}", schema)
	_, err = db.Exec(create_files_cmd)
	if err != nil {
		return err
	}

	create_raster_geoms_cmd := strings.ReplaceAll(sql_create_raster_geoms_table, "{{SCHEMA}}", schema)
	_, err = db.Exec(create_raster_geoms_cmd)
	if err != nil {
		return err
	}

	var projectID int
	query := `
		INSERT INTO public.projects (name)
		VALUES ($1)
		RETURNING project_id;
	`

	err = db.QueryRow(query, name).Scan(&projectID)
	if err != nil {
		return err
	}

	/*_, err = db.Exec("INSERT INTO public.projects (name) VALUES ($1);", name)
	if err != nil {
		return nil, err
	}*/

	slog.Debug("insert into table user_projects")

	_, err = db.Exec("INSERT INTO public.user_projects (user_id, project_id, role) VALUES ($1, $2, $3)", ownerid, projectID, "admin")
	return nil

}

/*
func (store *PSQLProjectStore) GetProject(id string) (*Project, error) { return nil, nil }
func (store *PSQLProjectStore) ListProjects() ([]*Project, error)      { return nil, nil }

func (store *PSQLProjectStore) DeleteProjectByName(name string) error { return nil }*/

// Create adds a new project to the database
/*
func (store *PSQLProjectStore) Create(project *Project) error {
	// Implementation for creating a project in the database

	db := store.db

	schema, err := psql.SanitizeSchemaName("project_" + project.Name)
	if err != nil {
		return err
	}

	slog.Debug("Adding project to table", "name", project.Name)
	_, err = db.Exec("INSERT INTO public.projects (name) VALUES ($1);", project.Name)
	if err != nil {
		return err
	}

	//create_tbl_cmd = strings.ReplaceAll(create_tbl_cmd, "{SCHEMA}", schema)
	create_schema_cmd := strings.ReplaceAll(sql_create_schema, "{{SCHEMA}}", schema)
	fmt.Println("Create schema: ", create_schema_cmd)
	_, err = db.Exec(create_schema_cmd)
	if err != nil {
		return err
	}

	create_files_cmd := strings.ReplaceAll(sql_create_files_table, "{{SCHEMA}}", schema)
	_, err = db.Exec(create_files_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Finished creating")

	create_raster_geoms_cmd := strings.ReplaceAll(sql_create_raster_geoms_table, "{{SCHEMA}}", schema)
	fmt.Println("Create table: ", create_raster_geoms_cmd)
	_, err = db.Exec(create_raster_geoms_cmd)
	if err != nil {
		return err
	}
	return nil
}*/

// GetByID retrieves a project by its ID from the database
func (store *PSQLProjectStore) GetByID(id string) (*Project, error) {
	// Implementation for retrieving a project by ID from the database
	return nil, errors.New("not implemented")
}

func (store *PSQLProjectStore) GetProjectsOwnedBy(userIDs []int) ([]*Project, error) {
	query := `
	SELECT 
	    p.project_id,
	    p.name,
	    p.created_at,
	    p.updated_at
	FROM 
	    public.projects p
	JOIN 
	    public.user_projects up ON p.project_id = up.project_id
	WHERE 
	    up.user_id = ANY($1); 
	`
	//var projects []*Project
	projects := []*Project{}
	if err := sqlscan.Select(context.Background(), store.db, &projects, query, pq.Array(userIDs)); err != nil {
		return nil, err
	}
	return projects, nil
}

// List retrieves all projects from the database
func (store *PSQLProjectStore) List() ([]*Project, error) {
	//cmd := "select regexp_replace(n1.schema_name, '^project_', '') as project from (select schema_name from information_schema.schemata where schema_name ~ '^project*') n1;"
	//rows, err := store.db.Query(sql_select_projects)
	rows, err := store.db.Query("SELECT name, created_at, updated_at FROM public.projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*Project{}
	for rows.Next() {
		var p Project
		err = rows.Scan(&p.Name, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
	// Implementation for listing all projects from the database
}

func (store *PSQLProjectStore) Delete(name string) error {

	schema, err := psql.SanitizeSchemaName(name)
	if err != nil {
		return err
	}

	cmd := strings.ReplaceAll(sql_drop_schema, "{{NAME}}", "project_"+schema)
	log.Println("Command: ", cmd)
	_, err = store.db.Exec(cmd)
	return err
}
