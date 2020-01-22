package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"lenslocked.com/views"
)

type Users struct {
	NewView *views.View
}

func (u *Users) New(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap-4", "views/users/new.gohtml"),
	}
}
