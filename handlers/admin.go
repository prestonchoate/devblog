package handlers

import (
	"log"
	"net/http"

	"github.com/prestonchoate/devblog/config"
)

func HandleAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	c := config.GetInstance()
	data := map[string]interface{}{
		"title":    "Admin Login",
		"links":    c.GetLinks(),
	}
	err := c.GetTemplates().ExecuteTemplate(w, "admin_login", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HandleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin dashboard"))
}

func AdminCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value != "1234" {
			log.Println("Session cookie not found")
			http.Redirect(w, r, "/dashboard/login", http.StatusTemporaryRedirect)
			return
		}
		log.Println("Session cookie validated")
		next.ServeHTTP(w, r)
	})
}
