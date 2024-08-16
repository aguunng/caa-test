package room

import "net/http"

type roomError struct {
	code int
}

const (
	roomErrorNotFound = iota
)

func (e *roomError) Error() string {
	switch e.code {
	case roomErrorNotFound:
		return "Room not found"
	default:
		return "Unknown error code"
	}
}

func (e *roomError) HTTPStatusCode() int {
	switch e.code {
	case roomErrorNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
