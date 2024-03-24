package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/prestonchoate/devblog/config"
)

type dbSessionPersister struct {
	table string
	keyName string
	dbConn *sql.DB
}

func NewDbSessionPersister(table string, keyname string) (*dbSessionPersister, error) {
	sessionPersister := dbSessionPersister{table: table, keyName: keyname}
	dbConfig := config.GetInstance().GetDBConfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DatabaseName))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	sessionPersister.dbConn = db

	err = sessionPersister.setup()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &sessionPersister, err
}

func (s *dbSessionPersister) setup() error {
	_, err := s.dbConn.Exec(
			fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT AUTO_INCREMENT PRIMARY KEY, session_id BINARY(16), created_at DATETIME, valid_until DATETIME",
		s.table,
		s.keyName),
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *dbSessionPersister) tableName() string {
	return s.table
}

func (s *dbSessionPersister) primaryKey() string {
	return s.keyName
}

func (s *dbSessionPersister) Save(session *Session) error {
	return fmt.Errorf("not implemented")
}

func (s *dbSessionPersister) SaveMany(session []*Session) error {
	return fmt.Errorf("not implemented")
}

func (s *dbSessionPersister) Load(id int) (*Session, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d", s.table, s.keyName, id)
	results, err := s.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	parsdResults := s.parseResults(results)
	if parsdResults == nil {
		return nil, fmt.Errorf("failed to parse results")
	}

	return parsdResults[0], nil
}

func (s *dbSessionPersister) LoadAll() ([]*Session, error) {
	query := fmt.Sprintf("SELECT * FROM %s", s.table)
	results, err := s.dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	parsedResults := s.parseResults(results)
	if parsedResults == nil {
		return nil, fmt.Errorf("failed to parse results")
	}

	return parsedResults, nil
}

func (s *dbSessionPersister) FilterBy(field string, val any) ([]*Session, error) {

	return nil, fmt.Errorf("not implemented")
}

func (d *dbSessionPersister) parseResults(results *sql.Rows) []*Session {
	sessions := []*Session{}
	for results.Next() {
		var id int
		var session_id uuid.UUID
		var user_id int
		var created_at time.Time
		var valid_until time.Time

		err := results.Scan(&id, &session_id, &user_id, &created_at, &valid_until)
		if err != nil {
			log.Println(err)
			return nil
		}

		sessions = append(sessions, &Session{
			Id: id,
			SessionId: session_id,
			UserId: user_id,
			CreatedAt: created_at,
			ValidUntil: valid_until,
		})
	}

	return sessions
}
