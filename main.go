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

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	router := mux.NewRouter()
	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.Handle("/faq", staticC.FAQ).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")

	log.Println("Server listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
