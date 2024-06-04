package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/health", app.Health)

	// Category
	mux.Post("/categories", app.CreateCategory)
	mux.Get("/categories/get/{id}", app.GetCategory)
	mux.Put("/categories/update/{id}", app.UpdateCategory)
	mux.Get("/categories/all", app.AllCategories)
	mux.Delete("/categories/delete/{id}", app.DeleteCategory)

	return mux
}
