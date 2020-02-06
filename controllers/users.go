package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
	Age      int8   `schema:"age"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap-4", "users/new"),
		us:      us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	must(u.NewView.Render(w, nil))
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	must(parseForm(r, &form))
	user := models.User{
		Name:  form.Name,
		Email: form.Email,
		Age:   form.Age,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}
