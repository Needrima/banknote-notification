package eto

type Event struct {
	EventReference     string
	EventName          string
	EventDate          string
	EventType          string
	EventSource        string
	EventUserReference string
	EventData          interface{}
}
