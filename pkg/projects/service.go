package projects

/*

// service implements the Service interface
type service struct {
	store ProjectStore
}

// NewService creates a new project service
func NewService(store ProjectStore) Service {
	return &service{
		store: store,
	}
}

// CreateProject creates a new project
func (s *service) CreateProject(name, description string, ownerid uint32) (*Project, error) {
	if name == "" || description == "" {
		return nil, errors.New("invalid project data")
	}

	project := &Project{
		//ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.store.Create(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject retrieves a project by its ID
func (s *service) GetProject(id string) (*Project, error) {
	if id == "" {
		return nil, errors.New("invalid project ID")
	}

	project, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) DeleteProjectByName(name string) error {
	return s.store.Delete(name)
}

// ListProjects lists all projects
func (s *service) ListProjects() ([]*Project, error) {
	return s.store.List()
}
*/
