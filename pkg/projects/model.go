package projects

import (
	"time"
)

// ProjectStore represents the project service interface
type ProjectStore interface {
	CreateProject(name, description string, ownerid int) error
	//GetProject(id string) (*Project, error)
	GetProjectsOwnedBy(userID []int) ([]*Project, error)
	//ListProjects() ([]*Project, error)
	//DeleteProjectByName(name string) error
}

type Project struct {
	WMSurl      string    `json:"wms,omitempty"`
	ID          int       `db:"project_id" json:"project_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ProductStore represents the interface for project storage operations
/*type ProjectStore interface {
	Create(project *Project) error
	GetByID(id string) (*Project, error)
	List() ([]*Project, error)
	Delete(name string) error
}*/
