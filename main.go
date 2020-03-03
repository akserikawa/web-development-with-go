package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
	"lenslocked.com/models"
	"lenslocked.com/rand"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "lenslocked.com"
)

func main() {
	config := DefaultConfig()
	dbConfig := DefaultPostgresConfig()
	services, err := models.NewServices(dbConfig.Dialect(), dbConfig.ConnectionInfo())
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery, services.Image, r)

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.Handle("/faq", staticController.FAQ).Methods("GET")
	r.Handle("/signup", usersController.NewView).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")
	r.Handle("/login", usersController.LoginView).Methods("GET")
	r.HandleFunc("/login", usersController.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	r.Handle("/galleries/new",
		requireUserMw.Apply(galleriesController.New)).
		Methods("GET")

	r.HandleFunc("/galleries",
		requireUserMw.ApplyFn(galleriesController.Create)).
		Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}",
		galleriesController.Show).
		Methods("GET").
		Name(controllers.ShowGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(galleriesController.Update)).
		Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(galleriesController.Delete)).
		Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(galleriesController.Edit)).
		Methods("GET").
		Name(controllers.EditGallery)

	r.Handle("/galleries",
		requireUserMw.ApplyFn(galleriesController.Index)).
		Methods("GET").
		Name(controllers.IndexGalleries)

	r.HandleFunc("/galleries/{id:[0-9]+}/images",
		requireUserMw.ApplyFn(galleriesController.ImageUpload)).
		Methods("POST")

	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesController.ImageDelete)).
		Methods("POST")

	assetHandler := http.FileServer(http.Dir("./assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetHandler))

	b, err := rand.Bytes(rand.RememberTokenBytes)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(config.IsProd()))

	log.Printf("Starting the server on http://localhost:%d ...\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port),
		csrfMw(userMw.Apply(r))))
}
