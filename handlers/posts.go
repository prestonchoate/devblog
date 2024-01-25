package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
)

func HandlePostListing(w http.ResponseWriter, r *http.Request) {
	c := config.GetInstance()
	postRepo, err := data.GetPostRepositoryInstance()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"title":    "Blog",
		"posts":    postRepo.GetPosts(0),
		"links":    c.GetLinks(),
		"showPost": true,
	}
	err = c.GetTemplates().ExecuteTemplate(w, "blog_listing", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIdParam := chi.URLParam(r, "postId")
		if postIdParam == "" {
			log.Println("Post ID not found")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(404), 404)
			return
		}
		postRepo, err := data.GetPostRepositoryInstance()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(404), 404)
			return
		}
		post := postRepo.GetPost(postId)
		if post == nil {
			log.Println("Post not found")
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "blog_post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HandlePostView(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value("blog_post").(*data.Post)
	if post == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	config := config.GetInstance()
	data := map[string]interface{}{
		"title": post.Title,
		"post":  post,
		"links": config.GetLinks(),
	}
	err := config.GetTemplates().ExecuteTemplate(w, "blog_post", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
