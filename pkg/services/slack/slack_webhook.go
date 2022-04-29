package slack

func parseWebhookResponse(raw []byte, v interface{}) error {
	var res = v.(**string)
	s := string(raw)
	*res = &s
	return nil
}
