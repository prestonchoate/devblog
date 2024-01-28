package data

type Persister[dataModel any] interface {
	tableName() string
	primaryKey() string
	setup() error
	Save(*dataModel) error
	SaveMany([]*dataModel) error
	Load(id int) (*dataModel, error)
	LoadAll() ([]*dataModel, error)
	FilterBy(field string, val any) ([]*dataModel, error)
}

// TODO: use this as an abstract way to create persister instances for any struct

type InMemoryPersister[dataModel any] struct {
	data   []*dataModel
	fields []string
}
