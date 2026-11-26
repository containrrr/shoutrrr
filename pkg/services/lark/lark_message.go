package lark

type RequestBody struct {
	MsgType   MsgType `json:"msg_type"`
	Content   Content `json:"content"`
	Timestamp string  `json:"timestamp,omitempty"`
	Sign      string  `json:"sign,omitempty"`
}

type MsgType string

const (
	MsgTypeText MsgType = "text"
	MsgTypePost MsgType = "post"
)

type Content struct {
	Text string `json:"text,omitempty"`
	Post *Post  `json:"post,omitempty"`
}

type Post struct {
	Zh *Message `json:"zh_cn,omitempty"`
	En *Message `json:"en_us,omitempty"`
}

type Message struct {
	Title   string   `json:"title"`
	Content [][]Item `json:"content"`
}

type Item struct {
	Tag  TagValue `json:"tag"`
	Text string   `json:"text,omitempty"`
	Link string   `json:"href,omitempty"`
}

type TagValue string

const (
	TagValueText TagValue = "text"
	TagValueLink TagValue = "a"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
