package main

import (
	"errors"
	"strings"
	"sync"
)

var events []Event
var registrations []Registration
var nextEventID = 1
var nextRegistrationID = 1
var mutex sync.Mutex

var errEventNotFound = errors.New("Мероприятие не найдено.")
var errDuplicateEmail = errors.New("Этот email уже зарегистрирован на выбранное мероприятие.")
var errNoSeats = errors.New("На мероприятии нет свободных мест.")

func getEvents() []EventInfo {
	mutex.Lock()
	defer mutex.Unlock()

	result := make([]EventInfo, 0, len(events))
	for _, event := range events {
		count := 0
		for _, registration := range registrations {
			if registration.EventID == event.ID {
				count++
			}
		}
		result = append(result, EventInfo{Event: event, Registered: count, Available: event.Capacity - count})
	}
	return result
}

func getEvent(id int) (EventInfo, []Registration, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, event := range events {
		if event.ID == id {
			var participants []Registration
			for _, registration := range registrations {
				if registration.EventID == id {
					participants = append(participants, registration)
				}
			}
			info := EventInfo{Event: event, Registered: len(participants), Available: event.Capacity - len(participants)}
			return info, participants, true
		}
	}
	return EventInfo{}, nil, false
}

func addEvent(title, date, location string, capacity int) {
	mutex.Lock()
	defer mutex.Unlock()

	event := Event{ID: nextEventID, Title: title, Date: date, Location: location, Capacity: capacity}
	nextEventID++
	events = append(events, event)
}

func addRegistration(eventID int, name, email string) error {
	mutex.Lock()
	defer mutex.Unlock()

	var event Event
	for _, item := range events {
		if item.ID == eventID {
			event = item
			break
		}
	}
	if event.ID == 0 {
		return errEventNotFound
	}

	email = strings.ToLower(strings.TrimSpace(email))
	count := 0
	for _, registration := range registrations {
		if registration.EventID == eventID {
			if registration.Email == email {
				return errDuplicateEmail
			}
			count++
		}
	}
	if count >= event.Capacity {
		return errNoSeats
	}

	registration := Registration{ID: nextRegistrationID, EventID: eventID, Name: name, Email: email}
	nextRegistrationID++
	registrations = append(registrations, registration)
	return nil
}

func cancelRegistration(id int) (int, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, registration := range registrations {
		if registration.ID == id {
			registrations = append(registrations[:i], registrations[i+1:]...)
			return registration.EventID, true
		}
	}
	return 0, false
}

func getStatistics() Statistics {
	mutex.Lock()
	defer mutex.Unlock()

	result := Statistics{EventCount: len(events), RegistrationCount: len(registrations)}
	for i, event := range events {
		count := 0
		for _, registration := range registrations {
			if registration.EventID == event.ID {
				count++
			}
		}
		if i == 0 || count > result.PopularCount {
			result.PopularEvent = event
			result.PopularCount = count
		}
	}
	return result
}
