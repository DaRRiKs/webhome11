package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func main() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", listEvents)
	mux.HandleFunc("GET /events/new", newEvent)
	mux.HandleFunc("POST /events", createEvent)
	mux.HandleFunc("GET /events/{id}", showEvent)
	mux.HandleFunc("GET /events/{id}/register", registrationForm)
	mux.HandleFunc("POST /events/{id}/register", registerParticipant)
	mux.HandleFunc("GET /events/{id}/participants", listParticipants)
	mux.HandleFunc("POST /registrations/{id}/cancel", cancelParticipant)
	mux.HandleFunc("GET /statistics", showStatistics)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
