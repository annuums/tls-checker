package slack

import "testing"

func TestGenerateMessage(t *testing.T) {

  txt := Text{
    Text: "Example",
    Bold: true,
  }

  msg, err := txt.GenerateMessage()
  if err != nil {

    t.Fatalf("unexpected error: %v", err)
  }

  if msg["type"] != "rich_text" {

    t.Errorf("expected type 'rich_text', got %v", msg["type"])
  }
}
