package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	attributes := map[string]string{
		"gender":  "Male",
		"country": "Japan",
	}

	data := struct {
		Name       string
		Age        int
		Attributes map[string]string
	}{"Akira", 25, attributes}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
