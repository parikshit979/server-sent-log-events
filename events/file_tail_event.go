package events

import (
	"fmt"
	"net/http"

	"github.com/server-sent-log-events/utils"
)

// FileTailEvent represents a file tail event.
type FileTailEvent struct {
	// source is the source of the event.
	source EventSource
	// eventType is the type of the event.
	eventType EventType
}

// NewFileTailEvent creates a new file tail event.
func NewFileTailEvent(source EventSource) *FileTailEvent {
	return &FileTailEvent{
		source:    source,
		eventType: EventTypeFileTail,
	}
}

// SendEventData returns the data of the event.
func (e *FileTailEvent) SendEventData(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fileUtil := utils.NewFileUtil("logfile.log", "log")
	err := fileUtil.OpenFile()
	if err != nil {
		fmt.Fprintf(w, "data: %s\n\n", err.Error())
		flusher.Flush()
	}

	// Create a channel to receive log lines.
	logChan := make(chan string)

	// Start tailing the file in a goroutine.
	go fileUtil.TailFile(r.Context(), logChan)
	defer fileUtil.CloseFile()

	// Write the log lines to the response.
	// This will block until the context is cancelled.
	for logLine := range logChan {
		// Send to SSE.
		fmt.Fprintf(w, "data: %s\n\n", logLine)
		flusher.Flush()
	}
}

// GetSource returns the source of the event.
func (e *FileTailEvent) GetSource() EventSource {
	return e.source
}

// GetType returns the type of the event.
func (e *FileTailEvent) GetType() EventType {
	return e.eventType
}
