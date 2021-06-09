# Webhook

Webhook is a bare bone package for working with discord webhooks.

## Example

simple plain text

```go
package main

import "github.com/Karitham/webhook"

const emojiURL = "https://cdn.discordapp.com/emojis/801874189526368256.gif"

func main() {
	wh := webhook.New("discord url here")

	wh.With(&webhook.Webhook{
		Username:  "Captain'Hook",
		AvatarURL: emojiURL,
		Content:   "This is the content of the message, it's plain text",
	})

	wh.Run()
}
```

with file as an attachement

```go
	resp, _ := http.Get(emojiURL)
	defer resp.Body.Close()

	wh.With(&webhook.Webhook{
		Username:  "nanabongo",
		AvatarURL: emojiURL,
		Files: []Attachment{{
				Body:     resp.Body,
				Filename: "nanabongo.gif",
			},
		},
	})

	wh.Run()
```

To have the file embed itself send an image with the embed object with it's url being `attachment://<filename>` and attach to file as an attachment
