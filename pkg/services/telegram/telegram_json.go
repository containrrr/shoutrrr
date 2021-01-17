package telegram

// SendMessagePayload is the notification payload for the telegram notification service
type SendMessagePayload struct {
	Text                string       `json:"text"`
	ID                  string       `json:"chat_id"`
	ParseMode           string       `json:"parse_mode,omitempty"`
	DisablePreview      bool         `json:"disable_web_page_preview"`
	DisableNotification bool         `json:"disable_notification"`
	ReplyMarkup         *ReplyMarkup `json:"reply_markup,omitempty"`
	Entities            []Entity     `json:"entities,omitempty"`
	ReplyTo             int64        `json:"reply_to_message_id"`
	MessageID           int64        `json:"message_id,omitempty"`
}

type UpdateMessagePayload struct {
	ChatID      string `json:"chat_id"`
	MessageID   int64  `json:"message_id"`
	Text        string `json:"text"`
	ParseMode   string `json:"parse_mode,omitempty"`
	ReplyMarkup string `json:"reply_markup,omitempty"`
}

type MessageResponse struct {
	OK     bool     `json:"ok"`
	Result *Message `json:"result"`
}

func createSendMessagePayload(message string, channel string, config *Config) SendMessagePayload {
	payload := SendMessagePayload{
		Text:                message,
		ID:                  channel,
		DisableNotification: !config.Notification,
		DisablePreview:      !config.Preview,
	}

	if config.ParseMode != ParseModes.None {
		payload.ParseMode = config.ParseMode.String()
	}

	return payload
}

type ErrorResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (e *ErrorResponse) Error() string {
	return e.Description
}

type UserResponse struct {
	OK     bool `json:"ok"`
	Result User `json:"result"`
}

type User struct {
	//	Unique identifier for this user or bot
	ID int64 `json:"id"`
	// True, if this user is a bot
	IsBot bool `json:"is_bot"`
	// User's or bot's first name
	FirstName string `json:"first_name"`
	//	Optional. User's or bot's last name
	LastName string `json:"last_name"`
	// Optional. User's or bot's username
	Username string `json:"username"`
	// Optional. IETF language tag of the user's language
	LanguageCode string `json:"language_code"`
	// 	Optional. True, if the bot can be invited to groups. Returned only in getMe.
	CanJoinGroups bool `json:"can_join_groups"`
	// 	Optional. True, if privacy mode is disabled for the bot. Returned only in getMe.
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages"`
	// 	Optional. True, if the bot supports inline queries. Returned only in getMe.
	SupportsInlineQueries bool `json:"supports_inline_queries"`
}

type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type CommandsRequest struct {
	Commands []Command `json:"commands"`
}

type UpdatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type UpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type InlineQuery struct {
	// Unique identifier for this query
	Id string `json:"id"`
	// Sender
	From User `json:"from"`
	// Text of the query (up to 256 characters)
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot
	Offset string `json:"offset"`
}

type InlineQueryResult struct {
	Type    string               `json:"result"`
	ID      string               `json:"id"`
	Title   string               `json:"title"`
	Content *InputMessageContent `json:"input_message_content"`
}

type InputMessageContent struct {
	Text string `json:"message_text"`
}

type InlineQueryAnswer struct {

	// Unique identifier for the answered query
	InlineQueryID string `json:"inline_query_id"`
	// A JSON-serialized array of results for the inline query
	Results []InlineQueryResult `json:"results"`
	// Optional	The maximum amount of time in seconds that the result of the inline query may be cached on the server. Defaults to 300.
	CacheTime int `json:"cache_time"`
	// Optional	Pass True, if results may be cached on the server side only for the user that sent the query. By default, results may be returned to any user who sends the same query
	IsPersonal bool `json:"is_personal"`
	// Optional	Pass the offset that a client should send in the next query with the same text to receive more results. Pass an empty string if there are no more results or if you don't support pagination. Offset length can't exceed 64 bytes.
	NextOffset string `json:"next_offset"`
	// Optional	If passed, clients will display a button with specified text that switches the user to a private chat with the bot and sends the bot a start message with the parameter switch_pm_parameter
	SwitchPMText string `json:"switch_pm_text"`
	// Optional	Deep-linking parameter for the /start message sent to the bot when user presses the switch button. 1-64 characters, only A-Z, a-z, 0-9, _ and - are allowed.
	SwitchPMParameter string `json:"switch_pm_parameter"`
}

type ChosenInlineResult struct {
}

type Update struct {
	// 	The update's unique identifier. Update identifiers start from a certain positive number and increase sequentially. This ID becomes especially handy if you're using Webhooks, since it allows you to ignore repeated updates or to restore the correct update sequence, should they get out of order. If there are no new updates for at least a week, then identifier of the next update will be chosen randomly instead of sequentially.
	UpdateID int `json:"update_id"`
	// 	Optional. New incoming message of any kind — text, photo, sticker, etc.
	Message *Message `json:"message"`
	// 	Optional. New version of a message that is known to the bot and was edited
	EditedMessage *Message `json:"edited_message"`
	// 	Optional. New incoming channel post of any kind — text, photo, sticker, etc.
	ChannelPost *Message `json:"channel_post"`
	// 	Optional. New version of a channel post that is known to the bot and was edited
	EditedChannelPost *Message `json:"edited_channel_post"`
	// 	Optional. New incoming inline query
	InlineQuery *InlineQuery `json:"inline_query"`
	//// 	Optional. The result of an inline query that was chosen by a user and sent to their chat partner. Please see our documentation on the feedback collecting for details on how to enable these updates for your bot.
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	//// 	Optional. New incoming callback query
	CallbackQuery *CallbackQuery `json:"callback_query"`
	//// 	Optional. New incoming shipping query. Only for invoices with flexible price
	//ShippingQuery	ShippingQuery `json:"shipping_query"`
	//// 	Optional. New incoming pre-checkout query. Contains full information about checkout
	//PreCheckoutQuery	PreCheckoutQuery `json:"pre_checkout_query"`
	/*
		// 	Optional. New poll state. Bots receive only updates about stopped polls and polls, which are sent by the bot
		Poll	Poll `json:"poll"`
		// 	Optional. A user changed their answer in a non-anonymous poll. Bots receive new votes only in polls that were sent by the bot itself.
		Poll_answer	PollAnswer `json:"poll_answer"`
	*/
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
}

func (c *Chat) Name() string {
	if c.Type == "private" || c.Type == "channel" {
		return "@" + c.Username
	}
	return c.Title
}

type InlineKey struct {
	Text                     string `json:"text"`
	Url                      string `json:"url"`
	LoginUrl                 string `json:"login_url"`
	CallbackData             string `json:"callback_data"`
	SwitchInlineQuery        string `json:"switch_inline_query"`
	SwitchInlineQueryCurrent string `json:"switch_inline_query_current_chat"`
}

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKey `json:"inline_keyboard,omitempty"`
}

type Entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"`
}

type CallbackQueryAnswer struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
}
