package telegram

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

// SendMessagePayload is the notification payload for the telegram notification service
type SendMessagePayload struct {
	Text                string       `json:"text"`
	ID                  string       `json:"chat_id"`
	MessageThreadID     *int         `json:"message_thread_id,omitempty"`
	ParseMode           string       `json:"parse_mode,omitempty"`
	DisablePreview      bool         `json:"disable_web_page_preview"`
	DisableNotification bool         `json:"disable_notification"`
	ReplyMarkup         *replyMarkup `json:"reply_markup,omitempty"`
	Entities            []entity     `json:"entities,omitempty"`
	ReplyTo             int64        `json:"reply_to_message_id"`
	MessageID           int64        `json:"message_id,omitempty"`
}

// Message represents one chat message
type Message struct {
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
}

type messageResponse struct {
	OK     bool     `json:"ok"`
	Result *Message `json:"result"`
}

func createSendMessagePayload(message string, channel string, config *Config) SendMessagePayload {
	var threadID *int = nil
	chatID, thread, ok := strings.Cut(channel, ":")
	if ok {
		if parsed, err := strconv.Atoi(thread); err == nil {
			threadID = &parsed
		}
	}
	payload := SendMessagePayload{
		Text:                message,
		ID:                  chatID,
		MessageThreadID:     threadID,
		DisableNotification: !config.Notification,
		DisablePreview:      !config.Preview,
	}

	parseMode := config.ParseMode
	if config.ParseMode == ParseModes.None && config.Title != "" {
		parseMode = ParseModes.HTML
		// no parse mode has been provided, treat message as unescaped HTML
		message = html.EscapeString(message)
	}

	if parseMode != ParseModes.None {
		payload.ParseMode = parseMode.String()
	}

	// only HTML parse mode is supported for titles
	if parseMode == ParseModes.HTML {
		payload.Text = fmt.Sprintf("<b>%v</b>\n%v", html.EscapeString(config.Title), message)
	}

	return payload
}

type errorResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func (e *errorResponse) Error() string {
	return e.Description
}

type userResponse struct {
	OK     bool `json:"ok"`
	Result User `json:"result"`
}

// User contains information about a telegram user or bot
type User struct {
	//	Unique identifier for this User or bot
	ID int64 `json:"id"`
	// True, if this User is a bot
	IsBot bool `json:"is_bot"`
	// User's or bot's first name
	FirstName string `json:"first_name"`
	//	Optional. User's or bot's last name
	LastName string `json:"last_name"`
	// Optional. User's or bot's username
	Username string `json:"username"`
	// Optional. IETF language tag of the User's language
	LanguageCode string `json:"language_code"`
	// 	Optional. True, if the bot can be invited to groups. Returned only in getMe.
	CanJoinGroups bool `json:"can_join_groups"`
	// 	Optional. True, if privacy mode is disabled for the bot. Returned only in getMe.
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages"`
	// 	Optional. True, if the bot supports inline queries. Returned only in getMe.
	SupportsInlineQueries bool `json:"supports_inline_queries"`
}

type updatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type updatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type inlineQuery struct {
	// Unique identifier for this query
	ID string `json:"id"`
	// Sender
	From User `json:"from"`
	// Text of the query (up to 256 characters)
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot
	Offset string `json:"offset"`
}

type chosenInlineResult struct{}

// Update contains state changes since the previous Update
type Update struct {
	// 	The Update's unique identifier. Update identifiers start from a certain positive number and increase sequentially. This ID becomes especially handy if you're using Webhooks, since it allows you to ignore repeated updates or to restore the correct Update sequence, should they get out of order. If there are no new updates for at least a week, then identifier of the next Update will be chosen randomly instead of sequentially.
	UpdateID int `json:"update_id"`
	// 	Optional. New incoming Message of any kind — text, photo, sticker, etc.
	Message *Message `json:"Message"`
	// 	Optional. New version of a Message that is known to the bot and was edited
	EditedMessage *Message `json:"edited_message"`
	// 	Optional. New incoming channel post of any kind — text, photo, sticker, etc.
	ChannelPost *Message `json:"channel_post"`
	// 	Optional. New version of a channel post that is known to the bot and was edited
	EditedChannelPost *Message `json:"edited_channel_post"`
	// 	Optional. New incoming inline query
	InlineQuery *inlineQuery `json:"inline_query"`
	//// 	Optional. The result of an inline query that was chosen by a User and sent to their chat partner. Please see our documentation on the feedback collecting for details on how to enable these updates for your bot.
	ChosenInlineResult *chosenInlineResult `json:"chosen_inline_result"`
	//// 	Optional. New incoming callback query
	CallbackQuery *callbackQuery `json:"callback_query"`

	// API fields that are not used by the client has been commented out

	//// 	Optional. New incoming shipping query. Only for invoices with flexible price
	//ShippingQuery	ShippingQuery `json:"shipping_query"`
	//// 	Optional. New incoming pre-checkout query. Contains full information about checkout
	//PreCheckoutQuery	PreCheckoutQuery `json:"pre_checkout_query"`
	//// 	Optional. New poll state. Bots receive only updates about stopped polls and polls, which are sent by the bot
	//Poll	Poll `json:"poll"`
	//// 	Optional. A User changed their answer in a non-anonymous poll. Bots receive new votes only in polls that were sent by the bot itself.
	//Poll_answer	PollAnswer `json:"poll_answer"`

	ChatMemberUpdate *ChatMemberUpdate `json:"my_chat_member"`
}

// Chat represents a telegram conversation
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
}

// Name returns the name of the channel based on its type
func (c *Chat) Name() string {
	if c.Type == "private" || c.Type == "channel" && c.Username != "" {
		return "@" + c.Username
	}
	return c.Title
}

type inlineKey struct {
	Text                     string `json:"text"`
	URL                      string `json:"url"`
	LoginURL                 string `json:"login_url"`
	CallbackData             string `json:"callback_data"`
	SwitchInlineQuery        string `json:"switch_inline_query"`
	SwitchInlineQueryCurrent string `json:"switch_inline_query_current_chat"`
}

type replyMarkup struct {
	InlineKeyboard [][]inlineKey `json:"inline_keyboard,omitempty"`
}

type entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type callbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"Message"`
	Data    string   `json:"data"`
}

// ChatMemberUpdate represents a member update in a telegram chat
type ChatMemberUpdate struct {
	// Chat the user belongs to
	Chat *Chat `json:"chat"`
	// Performer of the action, which resulted in the change
	From *User `json:"from"`
	// Date the change was done in Unix time
	Date int `json:"date"`
	// Previous information about the chat member
	OldChatMember *ChatMember `json:"old_chat_member"`
	// New information about the chat member
	NewChatMember *ChatMember `json:"new_chat_member"`
	// Optional. Chat invite link, which was used by the user to join the chat; for joining by invite link events only.
	// invite_link ChatInviteLink
}

// ChatMember represents the membership state for a user in a telegram chat
type ChatMember struct {
	//	The member's status in the chat
	Status string `json:"status"`
	// Information about the user
	User *User `json:"user"`
}
