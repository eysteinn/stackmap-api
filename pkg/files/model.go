package files

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          uuid.UUID `json:"uuid"`
	Filename    string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductStore represents the interface for project storage operations
type FileStore interface {
	Upload(project string, file *File, data []byte) error
	/*GetByID(project string, id uuid.UUID) (*File, error)
	List(project string) ([]*File, error)
	Delete(project string, id uuid.UUID) error*/
}
