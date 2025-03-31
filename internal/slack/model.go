package slack

type BlockMessage map[string]interface{}
type AttachmentMessage struct {
	Blocks   []BlockMessage `json:"blocks"`
	Fallback string         `json:"fallback,omitempty"`
}

type SlackMessage struct {
	ThreadTS    string              `json:"thread_ts,omitempty"` // for future feature
	Channel     string              `json:"channel"`
	Attachments []AttachmentMessage `json:"attachments,omitempty"`
	Blocks      []BlockMessage      `json:"blocks,omitempty"`
}

type SlackApiResponse struct {
	Channel          string `json:"channel"`
	Error            string `json:"error"`
	Ts               string `json:"ts"`
	OK               bool   `json:"ok"`
	ResponseMetadata struct {
		Warnings []string `json:"warnings"`
	} `json:"response_metadata"`
}
