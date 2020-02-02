package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	var id, orderID, orderAmount int
	var name, email, orderDescription string
	rows, err := db.Query(`
		SELECT users.id, users.name, users.email, orders.id, orders.amount, orders.description
		FROM users
		INNER JOIN orders
		ON users.id = orders.user_id`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		rows.Scan(&id, &name, &email, &orderID, &orderAmount, &orderDescription)
		fmt.Println("ID:", id, "Name:", name, "Email:", email,
			"Order ID:", orderID, "Amount:", orderAmount, "Description:", orderDescription)
	}

	db.Close()
}
