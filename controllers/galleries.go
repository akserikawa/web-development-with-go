package controllers

import (
	"net/http"

	"lenslocked.com/views"
)

type Galleries struct {
	NewView *views.View
}

func NewGalleries() *Galleries {
	return &Galleries{
		NewView: views.NewView("bootstrap-4", "galleries/new"),
	}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	must(g.NewView.Render(w, nil))
}
