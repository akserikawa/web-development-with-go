package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(" The page you requested could not be found."))
}

func main() {
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
