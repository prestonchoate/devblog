package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Post struct {
	Id          int
	Title       string
	Image       string
	Content     string
	Description string
}

type Project struct {
	ID          int
	Title       string
	Image       string
	Description string
	Link        string
}

var templates *template.Template
var templateDirs = []string{
	"templates/partials/*.html",
	"templates/*.html",
}

var posts = []Post{
	{Id: 1, Title: "Post 1", Image: "/public/placeholder.svg", Content: "Content 1", Description: "Description 1"},
	{Id: 2, Title: "Post 2", Image: "/public/placeholder.svg", Content: "Content 2", Description: "Description 2"},
	{Id: 3, Title: "Post 3", Image: "/public/placeholder.svg", Content: "Content 3", Description: "Description 3"},
	{Id: 4, Title: "Post 4", Image: "/public/placeholder.svg", Content: "Content 4", Description: "Description 4"},
	{Id: 5, Title: "Post 5", Image: "/public/placeholder.svg", Content: "Content 5", Description: "Description 5"},
	{Id: 6, Title: "Post 6", Image: "/public/placeholder.svg", Content: "Content 6", Description: "Description 6"},
	{Id: 7, Title: "Post 7", Image: "/public/placeholder.svg", Content: "Content 7", Description: "Description 7"},
}

var projects = []Project{
	{ID: 1, Title: "Project 1", Image: "/public/placeholder.svg", Description: "Description 1", Link: "https://github.com/prestonchoate"},
	{ID: 2, Title: "Project 2", Image: "/public/placeholder.svg", Description: "Description 2", Link: "https://github.com/prestonchoate"},
	{ID: 3, Title: "Project 3", Image: "/public/placeholder.svg", Description: "Description 3", Link: "https://github.com/prestonchoate"},
	{ID: 4, Title: "Project 4", Image: "/public/placeholder.svg", Description: "Description 4", Link: "https://github.com/prestonchoate"},
	{ID: 5, Title: "Project 5", Image: "/public/placeholder.svg", Description: "Description 5", Link: "https://github.com/prestonchoate"},
}

func loadTemplates() (templates *template.Template) {
	for _, dir := range templateDirs {
		var err error
		templates, err = templates.ParseGlob(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(templates.DefinedTemplates())
	return templates
}

func init() {
	templates = loadTemplates()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("homeHandler")
	// only send the first 3 posts and projects to the home page
	data := map[string]interface{}{
		"title":    "Home",
		"posts":    posts[:3],
		"projects": projects[:3],
	}
	err := templates.ExecuteTemplate(w, "home_page", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("listPosts")
	data := map[string]interface{}{
		"title": "Blog",
		"posts": posts,
	}
	err := templates.ExecuteTemplate(w, "blog_listing", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := chi.URLParam(r, "postId")
		var post *Post = nil
		for _, p := range posts {
			if strconv.Itoa(p.Id) == postId {
				post = &p
				break
			}
		}
		if post == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "blog_post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	fs := http.FileServer(http.Dir("./public/"))
	r.Mount("/public/", http.StripPrefix("/public/", fs))
	r.Get("/", homeHandler)

	r.Route("/blog", func(r chi.Router) {
		r.Get("/", listPosts)
		r.Route("/{postId}", func(r chi.Router) {
			r.Use(PostCtx)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				post := r.Context().Value("blog_post").(*Post)
				data := map[string]interface{}{
					"title": post.Title,
					"post":  post,
				}
				err := templates.ExecuteTemplate(w, "blog_post", data)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			})
		})

	})

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", r))

	//log.Fatal(http.ListenAndServe(":8080", r))
}
