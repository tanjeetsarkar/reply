package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/*", http.StripPrefix("/", fs))
	r.HandleFunc("/ws", serveWs)

	http.ListenAndServe(":5000", r)

}
