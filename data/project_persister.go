package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prestonchoate/devblog/config"
)

type memoryProjectPersister struct {
	projects []*Project
}

type dbProjectPersister struct {
	table   string
	keyName string
	dbConn  *sql.DB
}

func NewMemoryProjectPersister() (*memoryProjectPersister, error) {
	m := &memoryProjectPersister{}
	err := m.setup()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *memoryProjectPersister) LoadAll() ([]*Project, error) {
	return m.projects, nil
}

func (m *memoryProjectPersister) Load(id int) (*Project, error) {
	for _, project := range m.projects {
		if project.Id == id {
			return project, nil
		}
	}
	return nil, fmt.Errorf("Project with id %d not found", id)
}

func (m *memoryProjectPersister) Save(project *Project) error {
	m.projects = append(m.projects, project)
	return nil
}

func (m *memoryProjectPersister) SaveMany(projects []*Project) error {
	m.projects = append(m.projects, projects...)
	return nil
}

func (m *memoryProjectPersister) tableName() string {
	return "projects"
}

func (m *memoryProjectPersister) primaryKey() string {
	return "id"
}

func (m *memoryProjectPersister) setup() error {
	m.projects = []*Project{
		{Id: 1, Title: "Project 1", Image: "/public/placeholder.svg", Description: "Description 1", Link: "https://github.com/prestonchoate"},
		{Id: 2, Title: "Project 2", Image: "/public/placeholder.svg", Description: "Description 2", Link: "https://github.com/prestonchoate"},
		{Id: 3, Title: "Project 3", Image: "/public/placeholder.svg", Description: "Description 3", Link: "https://github.com/prestonchoate"},
		{Id: 4, Title: "Project 4", Image: "/public/placeholder.svg", Description: "Description 4", Link: "https://github.com/prestonchoate"},
		{Id: 5, Title: "Project 5", Image: "/public/placeholder.svg", Description: "Description 5", Link: "https://github.com/prestonchoate"},
	}
	return nil
}

func NewDBProjectPersister(table string, keyName string) (*dbProjectPersister, error) {
	d := &dbProjectPersister{table: table, keyName: keyName}
	dbConfig := config.GetInstance().GetDBConfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DatabaseName))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	d.dbConn = db
	err = d.setup()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return d, nil
}

func (m *dbProjectPersister) LoadAll() ([]*Project, error) {
	query := fmt.Sprintf("SELECT * FROM %s", m.table)
	results, err := m.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	parsedResults := m.parseResults(results)
	if parsedResults == nil {
		return nil, errors.New("failed to parse results")
	}

	return parsedResults, nil
}

func (m *dbProjectPersister) Load(id int) (*Project, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d", m.table, m.keyName, id)
	results, err := m.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	parsedResults := m.parseResults(results)
	if parsedResults == nil {
		return nil, errors.New("failed to parse results")
	}

	return parsedResults[0], nil
}

func (m *dbProjectPersister) Save(project *Project) error {
	return errors.New("not implemented")
}

func (m *dbProjectPersister) SaveMany(projects []*Project) error {
	return errors.New("not implemented")
}

func (m *dbProjectPersister) tableName() string {
	return m.table
}

func (m *dbProjectPersister) primaryKey() string {
	return m.keyName
}

func (m *dbProjectPersister) setup() error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s int NOT NULL AUTO_INCREMENT, title varchar(255), image varchar(255), description varchar(255), link varchar(255), PRIMARY KEY (id))", m.table, m.keyName)
	_, err := m.dbConn.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func (m *dbProjectPersister) parseResults(results *sql.Rows) []*Project {
	projects := []*Project{}
	for results.Next() {
		var project Project
		err := results.Scan(&project.Id, &project.Title, &project.Image, &project.Description, &project.Link)
		if err != nil {
			return nil
		}
		projects = append(projects, &project)
	}
	return projects
}
