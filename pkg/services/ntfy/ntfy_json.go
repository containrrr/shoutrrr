package ntfy

import "fmt"

type apiResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"error"`
	Link    string `json:"link"`
}

func (e *apiResponse) Error() string {
	msg := fmt.Sprintf("server response: %v (%v)", e.Message, e.Code)
	if e.Link != "" {
		return msg + ", see: " + e.Link
	}

	return msg
}
