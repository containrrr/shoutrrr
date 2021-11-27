package discourse

const postsEndpoint = "posts.json"

type createPostPayload struct {
	Title            string `json:"title"`
	Raw              string `json:"raw"`
	TopicId          *int   `json:"topic_id,omitempty"`
	Category         *int   `json:"category,omitempty"`
	TargetRecipients string `json:"target_recipients,omitempty"`
	Archetype        string `json:"archetype,omitempty"`
	CreatedAt        string `json:"created_at"`
	EmbedUrl         string `json:"embed_url"`
}

type createPostResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	CreatedAt    string `json:"created_at"`
	Raw          string `json:"raw"`
	Cooked       string `json:"cooked"`
	PostNumber   int    `json:"post_number"`
	PostType     int    `json:"post_type"`
	UpdatedAt    string `json:"updated_at"`
	Reads        int    `json:"reads"`
	ReadersCount int    `json:"readers_count"`
	Score        int    `json:"score"`
	Yours        bool   `json:"yours"`
	TopicId      int    `json:"topic_id"`
	TopicSlug    string `json:"topic_slug"`
	Version      int    `json:"version"`
}

type errorResponse struct {
	Action string   `json:"action"`
	Errors []string `json:"errors"`
}
