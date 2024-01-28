package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prestonchoate/devblog/config"
)

type dbUserPersister struct {
	table   string
	keyName string
	dbConn  *sql.DB
}

func NewDBUserPersister() (*dbUserPersister, error) {
	d := &dbUserPersister{
		table:   "users",
		keyName: "id",
	}
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

func (d *dbUserPersister) tableName() string {
	return d.table
}

func (d *dbUserPersister) primaryKey() string {
	return d.keyName
}

func (d *dbUserPersister) setup() error {
	createQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
		(%s INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255),
		passhash VARCHAR(255),
		UNIQUE(username)
		)`,
		d.table,
		d.keyName,
	)
	_, err := d.dbConn.Exec(createQuery)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *dbUserPersister) Save(user *User) error {
	query := fmt.Sprintf("INSERT INTO %s (username, passhash) VALUES (?, ?)", d.table)
	_, err := d.dbConn.Exec(query, user.Username, user.Passhash)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *dbUserPersister) SaveMany(users []*User) error {
	return fmt.Errorf("not implemented")
}

func (d *dbUserPersister) Load(id int) (*User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d", d.table, d.keyName, id)
	results, err := d.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	users, err := d.parseResults(results)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users[0], nil
}

func (d *dbUserPersister) LoadAll() ([]*User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", d.table)
	results, err := d.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	users, err := d.parseResults(results)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func (d *dbUserPersister) FilterBy(field string, val any) ([]*User, error) {
	query := fmt.Sprintf(`SELECT * From %s WHERE %s = %s`, d.tableName(), field, val)
	results, err := d.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	users, err := d.parseResults(results)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, err
}
func (d *dbUserPersister) parseResults(results *sql.Rows) ([]*User, error) {
	users := []*User{}
	for results.Next() {
		user := &User{}
		err := results.Scan(&user.Id, &user.Username, &user.Passhash)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

type memoryUserPersister struct {
	users []*User
}

func NewMemoryUserPersister() (*memoryUserPersister, error) {
	m := &memoryUserPersister{}
	err := m.setup()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *memoryUserPersister) LoadAll() ([]*User, error) {
	return m.users, nil
}

func (m *memoryUserPersister) Load(id int) (*User, error) {
	for _, user := range m.users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, nil
}

func (m *memoryUserPersister) Save(user *User) error {
	m.users = append(m.users, user)
	return nil
}

func (m *memoryUserPersister) SaveMany(users []*User) error {
	m.users = append(m.users, users...)
	return nil
}

func (m *memoryUserPersister) tableName() string {
	return "users"
}

func (m *memoryUserPersister) primaryKey() string {
	return "id"
}

func (m *memoryUserPersister) FilterBy(field string, val any) ([]*User, error) {
	fieldIndex := -1
	foundUsers := []*User{}
	for _, user := range m.users {
		if fieldIndex == -1 {
			fieldIndex = m.findFieldIndex(user, field)
			if fieldIndex == -1 {
				log.Println("Didn't find that field")
				return nil, errors.New("field does not exist on this type")
			}

			v := m.findFieldValue(user, fieldIndex)
			if v == val {
				log.Println("Found the user!")
				foundUsers = append(foundUsers, user)
			}
		}
	}

	if len(foundUsers) == 0 {
		log.Println("No user found")
		return nil, errors.New("User not found")
	}

	return foundUsers, nil
}

func (m *memoryUserPersister) findFieldValue(u *User, fieldIndex int) any {
	rv := reflect.Indirect(reflect.ValueOf(u))
	v := rv.Field(fieldIndex).Interface()
	return v
}

func (m *memoryUserPersister) findFieldIndex(u *User, field string) int {
	t := reflect.TypeOf(u)
	if t == nil {
		log.Println("Cannot filter this type of object")
		return -1
	}

	f, found := t.Elem().FieldByName(field)
	if !found {
		log.Println(field, " does not exist on this type")
		return -1
	}

	if len(f.Index) > 1 {
		log.Println("This field appears multiple times. Taking the first option")
	}
	if len(f.Index) <= 0 {
		return -1
	}

	return f.Index[0]
}

func (m *memoryUserPersister) setup() error {
	m.users = []*User{
		{
			Id:       1,
			Username: "test",
			Passhash: "password123",
		},
	}
	return nil
}
