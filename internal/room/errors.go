package room

import "net/http"

type roomError struct {
	code    int
	message string
}

const (
	roomErrorNotFound = iota
	agentNotFound     = iota
)

func (e *roomError) Error() string {
	if e.message != "" {
		return e.message
	}
	switch e.code {
	case roomErrorNotFound:
		return "Room not found"
	case agentNotFound:
		return "Agent not available"
	default:
		return "Unknown error code"
	}
}

func (e *roomError) HTTPStatusCode() int {
	switch e.code {
	case roomErrorNotFound, agentNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
