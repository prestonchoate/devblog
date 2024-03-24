package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/prestonchoate/devblog/config"
	"github.com/prestonchoate/devblog/data"
)

func HandleAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	c := config.GetInstance()
	data := map[string]interface{}{
		"title":    "Admin Login",
		"links":    c.GetLinks(),
		"messages": getErrorMessages(r),
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

func getErrorMessages(r *http.Request) []string {
	messageCookie, err := r.Cookie("messages")

	log.Printf("Cookie: %+v\n", messageCookie)

	if err == http.ErrNoCookie {
		return []string{}
	} else if err != nil {
		log.Println("Error getting messages cookie: ", err.Error())
		return []string{}
	}

	b64Str := messageCookie.Value
	decoded, err := base64.StdEncoding.DecodeString(b64Str)
	var messages data.ErrorMessageCookie
	err = json.Unmarshal(decoded, &messages)
	if err != nil {
		log.Println("Error unmarshalling messages cookie: ", err.Error())
		return []string{}
	}

	return messages.Messages
}
