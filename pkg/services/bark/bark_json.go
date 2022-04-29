package bark

// PushPayload is the notification payload for the bark notification service
type PushPayload struct {
	Body      string `json:"body"`
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Sound     string `json:"sound,omitempty"`
	Badge     *int64 `json:"badge,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	URL       string `json:"url,omitempty"`
	Category  string `json:"category,omitempty"`
	Copy      string `json:"copy,omitempty"`
}

type apiResponse struct {
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (e *apiResponse) Error() string {
	return "server response: " + e.Message
}
