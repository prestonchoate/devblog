package data

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prestonchoate/devblog/config"
)

type CmsPost struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Featured    bool   `json:"isFeatured"`
	Slug        string `json:"slug"`
	Image       string `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type hygraphPostsResponse struct {
	Data struct {
		Posts []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Content     struct {
				Html string `json:"html"`
			} `json:"content"`
			Featured bool   `json:"isFeatured"`
			Slug     string `json:"slug"`
			Image    struct {
				Url string `json:"url"`
			} `json:"image"`
			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"posts"`
	} `json:"data"`
}

func GetCmsPosts() []*CmsPost {
	//TODO: Utilize ValKey here to cache the posts instead of calling the API every time
	posts := make([]*CmsPost, 0, 40)
	lastId := ""

	cmsClient := config.GetInstance().GetCmsClient()
	response := getNextPostsPage(lastId, &cmsClient)

	if response == nil {
		return posts
	} else {
		for len(response.Data.Posts) > 0 {
			for _, postResponse := range response.Data.Posts {
				posts = append(posts, &CmsPost{
					ID:          postResponse.ID,
					Title:       postResponse.Title,
					Description: postResponse.Description,
					Content:     postResponse.Content.Html,
					Featured:    postResponse.Featured,
					Slug:        postResponse.Slug,
					Image:       postResponse.Image.Url,
					CreatedAt: postResponse.CreatedAt,
					UpdatedAt: postResponse.UpdatedAt,
				})
			}
			lastId = posts[len(posts)-1].ID
			response = getNextPostsPage(lastId, &cmsClient)
		}
	}

	return posts
}

func GetCmsPostBySlug(slug string) *CmsPost {
	req := fmt.Sprintf(`
			query {
				posts (where: {slug: "%v"}) {
					id
					title
					description
					content {
						html
					}
					isFeatured
					slug
					image {
						url
					}
					createdAt
					updatedAt
				}
			}
	`, slug)
	client := config.GetInstance().GetCmsClient()
	data := client.SendRequest(req)
	response := &hygraphPostsResponse{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	if len(response.Data.Posts) <= 0 {
		log.Println("No post found for Slug: ", slug)
		return nil
	}
	post := response.Data.Posts[0]
	return &CmsPost{
		ID: post.ID,
		Title: post.Title,
		Description: post.Description,
		Content: post.Content.Html,
		Featured: post.Featured,
		Slug: post.Slug,
		Image: post.Image.Url,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func getNextPostsPage(lastId string, client *config.HygraphClient) *hygraphPostsResponse {
	req := fmt.Sprintf(`query {
			posts (first: 100 %v) {
				id
				title
				description
				content {
					html
				}
				isFeatured
				slug
				image {
					url
				}
				createdAt
				updatedAt
			}
		}`, getPaginationField(lastId))

	data := client.SendRequest(req)
	response := &hygraphPostsResponse{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return response
}

func getPaginationField(lastId string) string {
	if lastId != "" {
		return fmt.Sprintf(", after: \"%v\"", lastId)
	} else {
		return ""
	}
}
