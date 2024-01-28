package data

import "errors"

var projectRepositoryInstance *ProjectRepository

type Project struct {
	Id          int
	Title       string
	Image       string
	Description string
	Link        string
}

type ProjectRepository struct {
	projects  []*Project
	persister Persister[Project]
}

func GetProjectRepositoryInstance() (*ProjectRepository, error) {
	if projectRepositoryInstance == nil {
		p, err := NewDBProjectPersister("projects", "id")
		if err != nil {
			return nil, errors.New("failed to create project repository")
		}
		projectRepositoryInstance = &ProjectRepository{}
		projectRepositoryInstance.persister = p
	}
	return projectRepositoryInstance, nil
}

func (pr *ProjectRepository) GetProjects(limit int) []*Project {
	projects, err := pr.persister.LoadAll()
	if err != nil {
		return nil
	}
	if limit <= 0 || limit >= len(projects) {
		return projects
	}
	return projects[:limit]
}

func (pr *ProjectRepository) GetProject(id int) *Project {
	project, err := pr.persister.Load(id)
	if err != nil {
		return nil
	}
	return project
}
