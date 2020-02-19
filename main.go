package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
	"lenslocked.com/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "lenslocked.com"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery, r)

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

	log.Println("Server listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", userMw.Apply(r)))
}
