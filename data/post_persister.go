package data

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prestonchoate/devblog/config"
)

type memoryPostPersister struct {
	posts []*Post
}

func randomDateTimeString() string {
	//generate a random date and return it as a string
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).String()
}

func NewMemoryPostPersister() (*memoryPostPersister, error) {
	m := &memoryPostPersister{}
	err := m.setup()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return m, nil
}

func (m *memoryPostPersister) LoadAll() ([]*Post, error) {
	return m.posts, nil
}

func (m *memoryPostPersister) Load(id int) (*Post, error) {
	for _, post := range m.posts {
		if post.Id == id {
			return post, nil
		}
	}
	return nil, fmt.Errorf("Post with id %d not found", id)
}

func (m *memoryPostPersister) Save(post *Post) error {
	m.posts = append(m.posts, post)
	return nil
}

func (m *memoryPostPersister) SaveMany(posts []*Post) error {
	m.posts = append(m.posts, posts...)
	return nil
}

func (m *memoryPostPersister) tableName() string {
	return "posts"
}

func (m *memoryPostPersister) primaryKey() string {
	return "id"
}

func (m *memoryPostPersister) setup() error {
	m.posts = []*Post{
		{
			Id:          1,
			Title:       "Post 1",
			Image:       "/public/placeholder.svg",
			Content:     "Content 1",
			Description: "Description 1",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
		{
			Id:          2,
			Title:       "Post 2",
			Image:       "/public/placeholder.svg",
			Content:     "Content 2",
			Description: "Description 2",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
		{
			Id:          3,
			Title:       "Post 3",
			Image:       "/public/placeholder.svg",
			Content:     "Content 3",
			Description: "Description 3",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
		{
			Id:          4,
			Title:       "Post 4",
			Image:       "/public/placeholder.svg",
			Content:     "Content 4",
			Description: "Description 4",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
		{
			Id:          5,
			Title:       "Post 5",
			Image:       "/public/placeholder.svg",
			Content:     "Content 5",
			Description: "Description 5",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
		{
			Id:          6,
			Title:       "Post 6",
			Image:       "/public/placeholder.svg",
			Content:     "Content 6",
			Description: "Description 6",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
		{
			Id:          7,
			Title:       "Post 7",
			Image:       "/public/placeholder.svg",
			Content:     "Content 7",
			Description: "Description 7",
			Featured:    false,
			CreatedAt:   randomDateTimeString(),
		},
	}
	return nil
}

type dbPostPersister struct {
	table   string
	keyName string
	dbConn  *sql.DB
}

func NewDBPostPersister(table string, keyName string) (*dbPostPersister, error) {
	d := dbPostPersister{table: table, keyName: keyName}
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
	return &d, err
}

func (d *dbPostPersister) setup() error {
	_, err := d.dbConn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255), image VARCHAR(255), content TEXT, description TEXT, featured BOOLEAN, created_at DATETIME)", d.table, d.keyName))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *dbPostPersister) tableName() string {
	return d.table
}

func (d *dbPostPersister) primaryKey() string {
	return d.keyName
}

func (d *dbPostPersister) Save(post *Post) error {
	return fmt.Errorf("not implemented")
}

func (d *dbPostPersister) SaveMany(posts []*Post) error {
	return fmt.Errorf("not implemented")
}

func (d *dbPostPersister) Load(id int) (*Post, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d", d.table, d.keyName, id)
	results, err := d.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	parsedResults := d.parseResults(results)
	if parsedResults == nil {
		return nil, fmt.Errorf("failed to parse results")
	}

	return parsedResults[0], nil
}

func (d *dbPostPersister) LoadAll() ([]*Post, error) {
	query := fmt.Sprintf("SELECT * FROM %s", d.table)
	results, err := d.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	parsedResults := d.parseResults(results)
	if parsedResults == nil {
		return nil, fmt.Errorf("failed to parse results")
	}

	return parsedResults, nil
}

func (d *dbPostPersister) parseResults(results *sql.Rows) []*Post {
	posts := []*Post{}
	for results.Next() {
		var id int
		var title string
		var image string
		var content string
		var description string
		var featured bool
		var createdAt string
		err := results.Scan(&id, &title, &image, &content, &description, &featured, &createdAt)
		if err != nil {
			log.Println(err)
			return nil
		}
		posts = append(posts, &Post{
			Id:          id,
			Title:       title,
			Image:       image,
			Content:     content,
			Description: description,
			Featured:    featured,
			CreatedAt:   createdAt,
		})
	}
	return posts
}
