package data

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var sessionRepositoryInstance *SessionRepository

type Session struct {
	Id int
	SessionId uuid.UUID
	UserId int
	CreatedAt time.Time
	ValidUntil time.Time
}

type SessionRepository struct {
	persister Persister[Session]
}

func GetSessionRepositoryInstance() (*SessionRepository, error) {
	if sessionRepositoryInstance == nil {
		p, err := NewDbSessionPersister("sessions", "id")
		if err != nil {
			return nil, errors.New("failed to create session repository")
		}
		sessionRepositoryInstance = &SessionRepository{
			persister: p,
		}
	}

	return sessionRepositoryInstance, nil
}

// TODO: implement creation and lookup of sessions
