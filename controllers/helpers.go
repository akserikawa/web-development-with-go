package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(request *http.Request, destination interface{}) error {
	if err := request.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	if err := decoder.Decode(destination, request.PostForm); err != nil {
		return err
	}
	return nil
}
