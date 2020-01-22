package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"lenslocked.com/views"
)

var homeView *views.View
var contactView *views.View
var faqView *views.View
var signupView *views.View

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func Contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func Faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(faqView.Render(w, nil))
}

func SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil))
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(" The page you requested could not be found."))
}

// A helper function that panics on any error
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("bootstrap-4", "views/home.gohtml")
	contactView = views.NewView("bootstrap-4", "views/contact.gohtml")
	faqView = views.NewView("bootstrap-4", "views/faq.gohtml")
	signupView = views.NewView("bootstrap-4", "views/signup.gohtml")

	router := httprouter.New()
	router.GET("/", Home)
	router.GET("/contact", Contact)
	router.GET("/faq", Faq)
	router.GET("/signup", SignUp)
	router.GET("/hello/:name", Hello)
	router.NotFound = http.HandlerFunc(NotFound)

	log.Println("Server listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
