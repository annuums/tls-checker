package slack

import "log"

type Text struct {
	Text string
	Bold bool
}

func (r Text) GenerateMessage() (map[string]interface{}, error) {

	style := map[string]bool{
		"bold": r.Bold,
	}

	return map[string]interface{}{
		"type": "rich_text",
		"elements": []map[string]interface{}{
			{
				"type": "rich_text_section",
				"elements": []map[string]interface{}{
					{
						"type":  "text",
						"text":  r.Text,
						"style": style,
					},
				},
			},
		},
	}, nil
}

type Payload struct {
	Channel     string `json:"channel"`
	Color       string `json:"color,omitempty"`
	Header      string `json:"header,omitempty"`
	Description string `json:"description,omitempty"`
	Text        []Text `json:"text"`
}

func (p *Payload) ToSlackMessage() *SlackMessage {

	if p.Color == "" {

		return p.convertToBlockMessage()
	}

	return p.convertToAttachmentMessage()
}

func (p *Payload) convertToAttachmentMessage() *SlackMessage {

	rawBlocks := p.generateRichText()
	blocks := make([]BlockMessage, len(rawBlocks))

	for i, item := range rawBlocks {

		blocks[i] = item
	}

	return &SlackMessage{
		Channel: p.Channel,
		Attachments: []AttachmentMessage{
			{
				Blocks:   blocks,
				Fallback: p.Description,
			},
		},
	}
}

func (p *Payload) convertToBlockMessage() *SlackMessage {

	rawBlocks := p.generateRichText()
	blocks := make([]BlockMessage, len(rawBlocks))

	for i, item := range rawBlocks {

		blocks[i] = item
	}

	return &SlackMessage{
		Channel: p.Channel,
		Blocks:  blocks,
	}
}

/*
-------------------------------------

--------------------------------------
*/
func (p *Payload) generateRichText() []map[string]interface{} {

	blocks := make([]map[string]interface{}, len(p.Text))

	for i, message := range p.Text {

		m, err := message.GenerateMessage()

		if err != nil {

			log.Printf("Error generating message for rich text %v: %v\n", message, err)

			text := message.Text

			blocks[i] = map[string]interface{}{
				"type": "rich_text",
				"elements": map[string]interface{}{
					"type": "rich_text_section",
					"elements": []map[string]interface{}{
						{
							"type": "text",
							"text": text,
							"style": map[string]bool{
								"bold":   false,
								"italic": false,
								"strike": false,
							},
						},
					},
				},
			}

			continue
		}
		blocks[i] = m
	}

	return blocks
}
