package data

var postRepositoryInstance *PostRepository

type Post struct {
	Id          int
	Title       string
	Image       string
	Content     string
	Description string
}

// TODO: replace posts in this with DB calls
type PostRepository struct {
	posts []*Post
}

func GetPostRepositoryInstance() *PostRepository {
	if postRepositoryInstance == nil {
		postRepositoryInstance = &PostRepository{
			posts: []*Post{
				{Id: 1, Title: "Post 1", Image: "/public/placeholder.svg", Content: "Content 1", Description: "Description 1"},
				{Id: 2, Title: "Post 2", Image: "/public/placeholder.svg", Content: "Content 2", Description: "Description 2"},
				{Id: 3, Title: "Post 3", Image: "/public/placeholder.svg", Content: "Content 3", Description: "Description 3"},
				{Id: 4, Title: "Post 4", Image: "/public/placeholder.svg", Content: "Content 4", Description: "Description 4"},
				{Id: 5, Title: "Post 5", Image: "/public/placeholder.svg", Content: "Content 5", Description: "Description 5"},
				{Id: 6, Title: "Post 6", Image: "/public/placeholder.svg", Content: "Content 6", Description: "Description 6"},
				{Id: 7, Title: "Post 7", Image: "/public/placeholder.svg", Content: "Content 7", Description: "Description 7"},
			},
		}
	}
	return postRepositoryInstance
}

func (pr *PostRepository) GetPosts(limit int) []*Post {
	if limit <= 0 {
		return pr.posts
	}

	return pr.posts[:limit]
}

func (pr *PostRepository) GetPost(id int) *Post {
	for _, post := range pr.posts {
		if post.Id == id {
			return post
		}
	}
	return nil
}
