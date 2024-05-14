package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
)

func HandlePostListing(w http.ResponseWriter, r *http.Request) {
	c := config.GetInstance()
	data := map[string]interface{}{
		"title":    "Blog",
		"posts":    data.GetCmsPosts(),
		"links":    c.GetLinks(),
		"showPost": true,
	}
	err := c.GetTemplates().ExecuteTemplate(w, "blog_listing", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slugParam := chi.URLParam(r, "slug")
		if slugParam == "" {
			log.Println("Slug not found")
			http.Error(w, http.StatusText(404), 404)
			return
		}
		post := data.GetCmsPostBySlug(slugParam)
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
	post := r.Context().Value("blog_post").(*data.CmsPost)
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
