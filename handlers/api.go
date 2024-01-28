package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/prestonchoate/devblog/data"
)

type userLoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (u *userLoginRequest) Bind(r *http.Request) error {
	log.Println("Binding user login request")
	log.Printf("userLoginRequest: %+v\n", u)
	return nil
}

func HandleApiLogin(w http.ResponseWriter, r *http.Request) {
	loginRequest := &userLoginRequest{}
	err := render.Bind(r, loginRequest)
	if err != nil {
		log.Printf("Error binding login request: \n%v\n", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ur, err := data.GetUserRepositoryInstance()
	if err != nil {
		log.Printf("Error getting user repository instance: \n%v\n", err.Error())
		http.Error(w, "Server Error", http.StatusBadGateway)
		return
	}

	log.Printf("Attempting to get user by username: %s\n", loginRequest.Username)
	user := ur.GetUserByUsername(loginRequest.Username)
	if user == nil {
		log.Println("Failed to find user by username: ", loginRequest.Username)
		http.Error(w, "Invalid Login Attempt", http.StatusBadRequest)
		return
	}

	if user.Passhash != loginRequest.Password {
		log.Println("Password incorrect")
		http.Error(w, "Invalid Login Attempt", http.StatusBadRequest)
		return
	}

	log.Println("Login successful. Creating session cookie")
	cookie := http.Cookie{
		Name:     "session",
		Value:    "1234",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/dashboard/admin", http.StatusSeeOther)

}
