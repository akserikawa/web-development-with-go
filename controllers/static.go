package controllers

import (
	"lenslocked.com/views"
)

type Static struct {
	Home    *views.View
	Contact *views.View
	FAQ     *views.View
}

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap-4", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap-4", "views/static/contact.gohtml"),
		FAQ:     views.NewView("bootstrap-4", "views/static/faq.gohtml"),
	}
}
