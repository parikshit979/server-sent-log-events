package events

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// ServerMonitoringEvent represents a server monitoring event.
type ServerMonitoringEvent struct {
	// source is the source of the event.
	source EventSource
	// eventType is the type of the event.
	eventType EventType
}

// NewServerMonitoringEvent creates a new server monitoring event.
func NewServerMonitoringEvent(source EventSource) *ServerMonitoringEvent {
	return &ServerMonitoringEvent{
		source:    source,
		eventType: EventTypeServerMonitoring,
	}
}

// SendEventData returns the data of the event.
func (e *ServerMonitoringEvent) SendEventData(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	monitoringTicker := time.NewTicker(time.Second)
	defer monitoringTicker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-monitoringTicker.C:
			memory, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("unable to get memory: %s", err.Error())
			}
			memUsage := fmt.Sprintf("Total: %d, Free:%d, UsedPercent:%.2f%%",
				memory.Total, memory.Free, memory.UsedPercent)

			c, err := cpu.Times(false)
			if err != nil {
				log.Printf("unable to get cpu: %s", err.Error())
			}
			cpuUsage := fmt.Sprintf("User: %f, Sys:%f, Idle:%f%%",
				c[0].User, c[0].System, c[0].Idle)

			monitoringData := fmt.Sprintf("Memory: [%s] ||| CPU: [%s]", memUsage, cpuUsage)
			fmt.Fprintf(w, "data: %s\n\n", monitoringData)
			flusher.Flush()
		}
	}

}

// GetSource returns the source of the event.
func (e *ServerMonitoringEvent) GetSource() EventSource {
	return e.source
}

// GetType returns the type of the event.
func (e *ServerMonitoringEvent) GetType() EventType {
	return e.eventType
}
