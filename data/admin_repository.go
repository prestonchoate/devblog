package data

var userRepositoryInstance *UserRepository

type User struct {
	Id       int
	Username string
	Passhash string
}

type UserRepository struct {
	persister Persister[User]
}

func GetUserRepositoryInstance() (*UserRepository, error) {
	if userRepositoryInstance == nil {
		//p, err := NewMemoryUserPersister()
		p, err := NewDBUserPersister()
		if err != nil {
			return nil, err
		}
		userRepositoryInstance = &UserRepository{
			persister: p,
		}
	}
	return userRepositoryInstance, nil
}

func (ur *UserRepository) GetUser(userId int) *User {
	u, err := ur.persister.Load(userId)
	if err != nil {
		return nil
	}
	return u
}

func (ur *UserRepository) GetUserByUsername(username string) *User {
	u, err := ur.persister.FilterBy("Username", username)
	if err != nil {
		return nil
	}
	// For now just return the first user found. There should really never be more than one
	return u[0]
}
