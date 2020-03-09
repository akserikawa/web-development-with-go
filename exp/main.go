package main

import (
	"fmt"

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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer services.Close()
	// services.DestructiveReset()
}
