package main

import (
	"net/http"
	"os"

	"github.com/server-sent-log-events/events"
	"github.com/server-sent-log-events/utils"
)

// sseFileTailHandler handles the server-sent events for file tailing.
// It streams log lines from a file to the client.
func sseFileTailHandler(w http.ResponseWriter, r *http.Request) {
	// // Create a new file tail event.
	fileTailEvent := events.NewFileTailEvent(events.EventSourceFile)

	// Send the event data to the client.
	fileTailEvent.SendEventData(w, r)
}

// sseServerMonitoringHandler handles the server-sent events for server monitoring.
// It streams server monitoring data to the client.
func sseServerMonitoringHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new server monitoring event.
	serverMonitoringEvent := events.NewServerMonitoringEvent(events.EventSourceServerMonitoring)

	// Send the event data to the client.
	serverMonitoringEvent.SendEventData(w, r)
}

func main() {
	// Simulate a log file for testing.
	go utils.SimulateLogFile("logfile.log", "log")

	// Set up the HTTP server and route.
	http.HandleFunc("/log/events", sseFileTailHandler)
	http.HandleFunc("/monitoring/events", sseServerMonitoringHandler)
	http.ListenAndServe(":8080", nil)

	os.Exit(0)
}
