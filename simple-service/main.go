package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	http.ListenAndServe(":42000", r)
}
