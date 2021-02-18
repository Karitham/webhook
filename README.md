# Webhook

Webhook is a bare bone package for working with discord webhooks.

## Example

```go
package main

import (
	"github.com/Karitham/webhook"
)

func main() {
	wh := webhook.New("discord url here")

	// Here we build the webhook object.
	// The only field required are either the content an Embed object
	wh.Webhook = &webhook.Webhook{
		Username:  "Captain'Hook",
		AvatarURL: "https://cdn.discordapp.com/emojis/801874189526368256.gif",
		Content:   "This is the content of the message, it's plain text",
	}

	if err := wh.Run(); err != nil {
		panic(err)
	}
}
```
