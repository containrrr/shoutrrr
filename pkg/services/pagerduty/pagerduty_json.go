package pagerduty

// EventPayload represents the payload being sent to the PagerDuty Events API v2
//
// See: https://developer.pagerduty.com/docs/events-api-v2-overview
type EventPayload struct {
	Payload     Payload `json:"payload"`
	RoutingKey  string  `json:"routing_key"`
	EventAction string  `json:"event_action"`
}

type Payload struct {
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}
