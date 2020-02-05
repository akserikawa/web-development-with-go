package main

import (
	"fmt"

	"lenslocked.com/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	us.DestructiveReset()

	user := models.User{
		Name:  "Akira Serikawa",
		Email: "akserikawa@gmail.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	foundUser, err := us.ByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)
}
