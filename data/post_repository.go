package data

import "errors"

var postRepositoryInstance *PostRepository

type Post struct {
	Id          int
	Title       string
	Image       string
	Content     string
	Description string
	Featured    bool
	CreatedAt   string
}

type PostRepository struct {
	persister Persister[Post]
}

func GetPostRepositoryInstance() (*PostRepository, error) {
	if postRepositoryInstance == nil {
		p, err := NewDBPostPersister("posts", "id")
		if err != nil {
			return nil, errors.New("failed to create post repository")
		}
		postRepositoryInstance = &PostRepository{
			persister: p,
		}
	}
	return postRepositoryInstance, nil
}

// TODO: Change this to add a limit to the persister instead of doing it here
func (pr *PostRepository) GetPosts(limit int) []*Post {
	p, err := pr.persister.LoadAll()
	if err != nil {
		return nil
	}
	if limit <= 0 || limit >= len(p) {
		return p
	}
	return p[:limit]
}

func (pr *PostRepository) GetPost(id int) *Post {
	p, err := pr.persister.Load(id)
	if err != nil {
		return nil
	}
	return p

}
