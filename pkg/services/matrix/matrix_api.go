package matrix

type MessageType string
type FlowType string

const (
	APILogin      = "/_matrix/client/r0/login"
	APIJoinInvite = "/_matrix/client/r0/rooms/%s/join"

	APISync        = "/_matrix/client/r0/sync"
	APISendMessage = "_matrix/client/r0/rooms/%s/send/m.room.message"
	APILookupRoom  = "/_matrix/client/r0/directory/room/%s"

	ContentType = "application/json"

	AccessTokenKey = "access_token"

	MsgTypeText       MessageType = "m.text"
	FlowLoginPassword FlowType    = "m.login.password"
)

type UserState struct {
	AccountData state  `json:"account_data"`
	NextBatch   string `json:"next_batch"`
	Presence    state  `json:"presence"`
	Rooms       rooms  `json:"rooms"`
}

type APIResLogin struct {
	Flows []flow `json:"flows"`
}

type APIReqLoginPassword struct {
	Type     FlowType `json:"type"`
	User     string   `json:"user"`
	Password string   `json:"password"`
}

type APIResLoginPassword struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
}

type APIReqSend struct {
	MsgType MessageType `json:"msgtype"`
	Body    string      `json:"body"`
}

type APIResRoom struct {
	RoomID string `json:"room_id"`
}

type APIResEvent struct {
	EventID string `json:"event_id"`
}

type flow struct {
	Type FlowType `json:"type"`
}

type rooms struct {
	Invite map[string]room `json:"invite"`
	Join   map[string]room `json:"join"`
}

type room struct {
	AccountData         state       `json:"account_data"`
	Ephemeral           state       `json:"ephemeral"`
	State               state       `json:"state"`
	Timeline            state       `json:"timeline"`
	UnreadNotifications interface{} `json:"unread_notifications"`
}

type state struct {
	Events    []event `json:"events"`
	Limited   bool    `json:"limited"`
	PrevBatch string  `json:"prev_batch"`
}

type event struct {
	Content        map[string]interface{} `json:"content"`
	EventID        string                 `json:"event_id"`
	OriginServerTS uint64                 `json:"origin_server_ts"`
	Sender         string                 `json:"sender"`
	StateKey       string                 `json:"state_key"`
	Type           string                 `json:"type"`
	Unsigned       unsigned               `json:"unassigned"`
}

type unsigned struct {
	Age uint `json:"age"`
}
