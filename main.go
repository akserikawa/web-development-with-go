package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
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

	userService, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer userService.Close()
	userService.AutoMigrate()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(userService)

	router := mux.NewRouter()
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.Handle("/faq", staticController.FAQ).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")
	router.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	log.Println("Server listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
