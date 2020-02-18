package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
	Age      int8   `schema:"age"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap-4", "users/new"),
		LoginView: views.NewView("bootstrap-4", "users/login"),
		us:        us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	must(u.NewView.Render(w, nil))
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Age:      form.Age,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
		return
	}

	refererCookie, err := r.Cookie("referer")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, refererCookie.Value, http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")

	refererCookie := http.Cookie{
		Name:     "referer",
		Value:    r.URL.EscapedPath(),
		HttpOnly: true,
	}

	if err != nil {
		http.SetCookie(w, &refererCookie)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.SetCookie(w, &refererCookie)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintln(w, user)
}
