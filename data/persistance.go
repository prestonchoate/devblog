package data

type Persister[dataModel any] interface {
	tableName() string
	primaryKey() string
	setup() error
	Save(*dataModel) error
	SaveMany([]*dataModel) error
	Load(id int) (*dataModel, error)
	LoadAll() ([]*dataModel, error)
}
