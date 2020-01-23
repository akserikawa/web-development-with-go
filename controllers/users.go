package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/julienschmidt/httprouter"
	"lenslocked.com/views"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap-4", "views/users/new.gohtml"),
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec := schema.NewDecoder()
	form := SignupForm{}
	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}
