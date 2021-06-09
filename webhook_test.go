package webhook

import (
	"net/http"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	const emojiURL = "https://cdn.discordapp.com/emojis/801874189526368256.gif"

	url := os.Getenv("WEBHOOK_URL")
	if url == "" {
		t.Skip("no webhook URL provided")
	}
	wh := New(url)

	webhooks := []*Webhook{
		{
			Username:  "Captain'Hook",
			AvatarURL: emojiURL,
			Content:   "This is the content of the message, it's plain text",
		},
		{
			Username:  "Captain'Hook",
			AvatarURL: emojiURL,
			Embeds: []Embed{{
				Description: "This is the body of the embed",
				Title:       "This is the title",
				Footer:      Footer{Text: "This is the footer", IconURL: emojiURL},
				Color:       0xB00B69,
				Thumbnail:   Image{URL: emojiURL},
				Author:      Author{Name: "The author", URL: emojiURL, IconURL: emojiURL},
				Image:       Image{URL: emojiURL},
				URL:         emojiURL,
				Fields: []Field{
					{Name: "1", Value: "nana", Inline: true},
					{Name: "2", Value: "bongo", Inline: true},
					{Name: "3", Value: "nanabongo", Inline: false},
				},
			}},
		},
	}

	resp, err := http.Get(emojiURL)
	if err == nil {
		webhooks = append(webhooks, &Webhook{
			Username:  "nanabongo",
			AvatarURL: emojiURL,
			Files: []Attachment{{
				Body:     resp.Body,
				Filename: "nanabongo.gif",
			}},
		})
		defer resp.Body.Close()
	}

	for _, w := range webhooks {
		wh.With(w)

		if err := wh.Run(); err != nil {
			t.Fatal(err)
		}
	}
}
