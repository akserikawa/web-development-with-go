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

	router := mux.NewRouter()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery, router)

	requireUserMiddleware := middleware.RequireUser{
		UserService: services.User,
	}

	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.Handle("/faq", staticController.FAQ).Methods("GET")
	router.Handle("/signup", usersController.NewView).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")
	router.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	newGallery := requireUserMiddleware.Apply(galleriesController.New)
	createGallery := requireUserMiddleware.ApplyFn(galleriesController.Create)
	editGallery := requireUserMiddleware.ApplyFn(galleriesController.Edit)
	updateGallery := requireUserMiddleware.ApplyFn(galleriesController.Update)
	deleteGallery := requireUserMiddleware.ApplyFn(galleriesController.Delete)

	router.Handle("/galleries/new", newGallery).Methods("GET")
	router.HandleFunc("/galleries", createGallery).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleriesController.Show).
		Methods("GET").Name(controllers.ShowGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET")
	router.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")

	log.Println("Server listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
