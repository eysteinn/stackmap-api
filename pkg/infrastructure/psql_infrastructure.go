package infrastructure

import (
	"database/sql"
	_ "embed"
	"log/slog"
)

//go:embed sql_create_projects_table.sql
var sql_create_projects_table string

//go:embed sql_create_user_table.sql
var sql_create_user_table string

//go:embed sql_create_user_project_table.sql
var sql_create_user_project_table string

//go:embed sql_delete_infrastructure.sql
var sql_delete_infrastructure string

//go:embed sql_create_refreshtokens_table.sql
var sql_create_refresh_tokens_table string

/*type Infrastructure interface {
	CreateInrastructure() error
}*/

// PSQLProjectStore implements ProjectStore interface using an SQL database
type psqlInfrastructure struct {
	db *sql.DB
}

func NewPSQLInfrastructure(db *sql.DB) *psqlInfrastructure {
	return &psqlInfrastructure{db: db}
}
func (infra *psqlInfrastructure) Delete() error {
	slog.Debug("Deleting infrastructure")
	_, err := infra.db.Exec(sql_delete_infrastructure)
	if err != nil {
		slog.Debug("Error while deleting infrastructure", "Error", err)
	}
	return err
}
func (infra *psqlInfrastructure) Create() (err error) {
	db := infra.db
	slog.Debug("Creating infrastructure", "table", "projects")
	if _, err = db.Exec(sql_create_projects_table); err != nil {
		return err
	}

	slog.Debug("Creating infrastructure", "table", "users")
	if _, err = db.Exec(sql_create_user_table); err != nil {
		return err
	}

	slog.Debug("Creating infrastructure", "table", "user_projects")
	if _, err = db.Exec(sql_create_user_project_table); err != nil {
		return err
	}

	slog.Debug("Creating infrastructur", "table", "refresh_tokens")
	if _, err = db.Exec(sql_create_refresh_tokens_table); err != nil {
		return err
	}

	return nil
}
