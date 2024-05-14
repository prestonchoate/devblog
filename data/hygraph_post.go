package data

import (
	"encoding/json"
	"fmt"
	"log"

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
		} `json:"posts"`
	} `json:"data"`
}

func GetCmsPosts() []*CmsPost {
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
				})
			}
			lastId = posts[len(posts)-1].ID
			response = getNextPostsPage(lastId, &cmsClient)
		}
	}

	return posts
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
