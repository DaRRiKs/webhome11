package main

import (
	"bytes"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type pageData struct {
	Title         string
	Error         string
	Form          url.Values
	Events        []EventInfo
	Event         EventInfo
	Registrations []Registration
	Statistics    Statistics
}

func render(w http.ResponseWriter, status int, name string, data pageData) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, name, data); err != nil {
		log.Println(err)
		http.Error(w, "Не удалось отобразить страницу.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := buffer.WriteTo(w); err != nil {
		log.Println(err)
	}
}

func readID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный идентификатор: нужно целое число больше нуля.", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "events.html", pageData{Title: "Мероприятия", Events: getEvents()})
}

func newEvent(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "new_event.html", pageData{Title: "Добавление мероприятия"})
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Не удалось прочитать данные формы.", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.PostForm.Get("title"))
	date := strings.TrimSpace(r.PostForm.Get("date"))
	location := strings.TrimSpace(r.PostForm.Get("location"))
	capacity, capacityErr := strconv.Atoi(strings.TrimSpace(r.PostForm.Get("capacity")))
	_, dateErr := time.Parse("2006-01-02", date)
	data := pageData{Title: "Добавление мероприятия", Form: r.PostForm}

	switch {
	case title == "":
		data.Error = "Укажите название мероприятия."
	case dateErr != nil:
		data.Error = "Укажите корректную дату в формате YYYY-MM-DD."
	case location == "":
		data.Error = "Укажите место проведения."
	case capacityErr != nil || capacity <= 0:
		data.Error = "Вместимость должна быть целым числом больше нуля."
	}
	if data.Error != "" {
		render(w, http.StatusBadRequest, "new_event.html", data)
		return
	}

	addEvent(title, date, location, capacity)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func showEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	event, _, found := getEvent(id)
	if !found {
		http.Error(w, errEventNotFound.Error(), http.StatusNotFound)
		return
	}
	render(w, http.StatusOK, "event.html", pageData{Title: event.Title, Event: event})
}

func registrationForm(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	event, _, found := getEvent(id)
	if !found {
		http.Error(w, errEventNotFound.Error(), http.StatusNotFound)
		return
	}
	render(w, http.StatusOK, "register.html", pageData{Title: "Регистрация участника", Event: event})
}

func registerParticipant(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	event, _, found := getEvent(id)
	if !found {
		http.Error(w, errEventNotFound.Error(), http.StatusNotFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Не удалось прочитать данные формы.", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.PostForm.Get("name"))
	email := strings.TrimSpace(r.PostForm.Get("email"))
	address, emailErr := mail.ParseAddress(email)
	data := pageData{Title: "Регистрация участника", Event: event, Form: r.PostForm}
	if name == "" {
		data.Error = "Укажите имя участника."
	} else if emailErr != nil || address.Address != email {
		data.Error = "Укажите корректный email, например name@example.com."
	}
	if data.Error != "" {
		render(w, http.StatusBadRequest, "register.html", data)
		return
	}

	if err := addRegistration(id, name, email); err != nil {
		if err == errEventNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		data.Error = err.Error()
		render(w, http.StatusConflict, "register.html", data)
		return
	}
	http.Redirect(w, r, "/events/"+strconv.Itoa(id)+"/participants", http.StatusSeeOther)
}

func listParticipants(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	event, registrations, found := getEvent(id)
	if !found {
		http.Error(w, errEventNotFound.Error(), http.StatusNotFound)
		return
	}
	render(w, http.StatusOK, "participants.html", pageData{Title: "Участники мероприятия", Event: event, Registrations: registrations})
}

func cancelParticipant(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	eventID, found := cancelRegistration(id)
	if !found {
		http.Error(w, "Регистрация не найдена.", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/events/"+strconv.Itoa(eventID)+"/participants", http.StatusSeeOther)
}

func showStatistics(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "statistics.html", pageData{Title: "Статистика", Statistics: getStatistics()})
}
