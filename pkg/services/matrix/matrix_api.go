package matrix

type MessageType string
type FlowType string
type IdentifierType string

const (
	APILogin       = "/_matrix/client/r0/login"
	APIRoomJoin    = "/_matrix/client/r0/join/%s"
	APISendMessage = "/_matrix/client/r0/rooms/%s/send/m.room.message"
	APIJoinedRooms = "/_matrix/client/r0/joined_rooms"

	ContentType = "application/json"

	AccessTokenKey = "access_token"

	MsgTypeText       MessageType    = "m.text"
	FlowLoginPassword FlowType       = "m.login.password"
	IDTypeUser        IdentifierType = "m.id.user"
)

type APIResLoginFlows struct {
	Flows []flow `json:"flows"`
}

type APIReqLogin struct {
	Type       FlowType    `json:"type"`
	Identifier *identifier `json:"identifier"`
	Password   string      `json:"password,omitempty"`
	Token      string      `json:"token,omitempty"`
}

type APIResLogin struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
}

type APIReqSend struct {
	MsgType MessageType `json:"msgtype"`
	Body    string      `json:"body"`
}

type APIResRoom struct {
	RoomID string `json:"room_id"`
}

type APIResJoinedRooms struct {
	Rooms []string `json:"joined_rooms"`
}

type APIResEvent struct {
	EventID string `json:"event_id"`
}

type APIResError struct {
	Message string `json:"error"`
	Code    string `json:"errcode"`
}

func (e *APIResError) Error() string {
	return e.Message
}

type flow struct {
	Type FlowType `json:"type"`
}

type identifier struct {
	Type IdentifierType `json:"type"`
	User string         `json:"user,omitempty"`
}

func NewUserIdentifier(user string) (id *identifier) {
	return &identifier{
		Type: IDTypeUser,
		User: user,
	}
}
