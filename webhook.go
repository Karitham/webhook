package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Hook is the toplevel object.
// It contains a Webhook object and everything required to send it.
// It's optimised for reusability, so it has an embedded http.Client.
// You can modify the client yourself if you want to change the defaults.
type Hook struct {
	Webhook *Webhook
	encoder *json.Encoder
	Client  *http.Client
	buf     *bytes.Buffer
	url     string
}

// Webhook is the webhook object sent to discord
type Webhook struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds"`
}

// Embed is the embed object
type Embed struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Color       int64   `json:"color"`
	Fields      []Field `json:"fields"`
	Thumbnail   Image   `json:"thumbnail"`
	Image       Image   `json:"image"`
	Footer      Footer  `json:"footer"`
}

// Author is the author object
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// Field is the field object inside an embed
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Footer is the footer of the embed
type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

// Image is an image possibly contained inside the embed
type Image struct {
	URL string `json:"url"`
}

// New returns a new webhook with the designated URL.
func New(URL string) *Hook {
	buffer := &bytes.Buffer{}
	return &Hook{
		Webhook: &Webhook{},
		url:     URL,
		Client:  &http.Client{Timeout: time.Second * 10},
		encoder: json.NewEncoder(buffer),
		buf:     buffer,
	}
}

// Run the webhook with the preconfigured settings
func (w *Hook) Run() error {
	// Encode the body
	err := w.encoder.Encode(w.Webhook)
	if err != nil {
		return fmt.Errorf("error encoding the webhook: %w", err)
	}

	// Building the request
	req, err := http.NewRequest(http.MethodPost, w.url, w.buf)
	if err != nil {
		return fmt.Errorf("error while building the request: %w", err)
	}
	req.Header.Add("content-type", "application/json")

	// Sending the rq
	resp, err := w.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending the webhook: %w", err)
	}
	defer resp.Body.Close()

	// Normal behavior should return here
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// Error handling and wrapping
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading the response body: %w. status code: %s", err, resp.Status)
	}
	return fmt.Errorf("error sending the webhook: %s, status code: %s", body, resp.Status)
}
