package main

import (
	"fmt"

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

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)

	db.AutoMigrate(&User{}, &Order{})

	var u User
	db.First(&u)
	if db.Error != nil {
		panic(db.Error)
	}

	createOrder(db, u, 100, "Castelli Bib Shorts")
	createOrder(db, u, 10, "Bus Ticket")
	createOrder(db, u, 29, "Summer Festival Ticket")
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}
