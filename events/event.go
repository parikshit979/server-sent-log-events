package events

import "net/http"

type Event interface {
	// SendEventData returns the data of the event.
	SendEventData(w http.ResponseWriter, r *http.Request)
	// GetSource returns the source of the event.
	GetSource() string
	// GetType returns the type of the event.
	GetType() string
}

// EventType is the type of the event.
type EventType string

const (
	// EventTypeFileTail is the type of file tail event.
	EventTypeFileTail EventType = "file_tail"
	// EventTypeSystem is the type of system event.
	EventTypeServerMonitoring EventType = "server_monitoring"
)

// EventSource is the source of the event.
type EventSource string

const (
	// EventSourceUser is the source of events from file tail.
	EventSourceFile EventSource = "file"
	// EventSourceServerMonitoring is the source of events from server monitoring.
	EventSourceServerMonitoring EventSource = "server_monitoring"
)
