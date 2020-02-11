package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"lenslocked.com/hash"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "lenslocked.com"
)

func main() {
	hmac := hash.NewHMAC("my-secret-key")

	fmt.Println(hmac.Hash("this is my string to hash"))
}
