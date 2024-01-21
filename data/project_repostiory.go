package data

var projectRepositoryInstance *ProjectRepository

type Project struct {
	Id          int
	Title       string
	Image       string
	Description string
	Link        string
}

type ProjectRepository struct {
	projects []*Project
}

func GetProjectRepositoryInstance() *ProjectRepository {
	if projectRepositoryInstance == nil {
		projectRepositoryInstance = &ProjectRepository{
			projects: []*Project{
				{Id: 1, Title: "Project 1", Image: "/public/placeholder.svg", Description: "Description 1", Link: "https://github.com/prestonchoate"},
				{Id: 2, Title: "Project 2", Image: "/public/placeholder.svg", Description: "Description 2", Link: "https://github.com/prestonchoate"},
				{Id: 3, Title: "Project 3", Image: "/public/placeholder.svg", Description: "Description 3", Link: "https://github.com/prestonchoate"},
				{Id: 4, Title: "Project 4", Image: "/public/placeholder.svg", Description: "Description 4", Link: "https://github.com/prestonchoate"},
				{Id: 5, Title: "Project 5", Image: "/public/placeholder.svg", Description: "Description 5", Link: "https://github.com/prestonchoate"},
			},
		}
	}
	return projectRepositoryInstance
}

func (pr *ProjectRepository) GetProjects(limit int) []*Project {
	if limit <= 0 {
		return pr.projects
	}
	return pr.projects[:limit]
}

func (pr *ProjectRepository) GetProject(id int) *Project {
	for _, project := range pr.projects {
		if project.Id == id {
			return project
		}
	}
	return nil
}
