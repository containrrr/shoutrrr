package gotify

import "fmt"

// messageRequest is the actual payload being sent to the Gotify API
type messageRequest struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type messageResponse struct {
	messageRequest
	ID    uint64 `json:"id"`
	AppID uint64 `json:"appid"`
	Date  string `json:"date"`
}

type errorResponse struct {
	Name        string `json:"error"`
	Code        uint64 `json:"errorCode"`
	Description string `json:"errorDescription"`
}

func (er *errorResponse) Error() string {
	return fmt.Sprintf("server respondend with %v (%v): %v", er.Name, er.Code, er.Description)
}
