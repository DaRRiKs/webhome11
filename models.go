package main

type Event struct {
	ID       int
	Title    string
	Date     string
	Location string
	Capacity int
}

type Registration struct {
	ID      int
	EventID int
	Name    string
	Email   string
}

type EventInfo struct {
	Event
	Registered int
	Available  int
}

type Statistics struct {
	EventCount        int
	RegistrationCount int
	PopularEvent      Event
	PopularCount      int
}
