package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"
)

// Hook is the toplevel object.
// It contains a Webhook object and everything required to send it.
// It's optimised for reusability, so it has an embedded http.Client.
// You can modify the client yourself if you want to change the defaults.
type Hook struct {
	Webhook *Webhook
	Client  *http.Client
	mutex   sync.Mutex
	url     string
}

// Webhook is the webhook object sent to discord
type Webhook struct {
	Username  string       `json:"username"`
	AvatarURL string       `json:"avatar_url"`
	Content   string       `json:"content"`
	Embeds    []Embed      `json:"embeds"`
	Files     []Attachment `json:"-"`
}

// Attachement is the files attached to the request
type Attachment struct {
	Body     io.Reader
	Filename string
}

// Embed is the embed object
type Embed struct {
	Author      Author  `json:"author"`
	Footer      Footer  `json:"footer"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Thumbnail   Image   `json:"thumbnail"`
	Image       Image   `json:"image"`
	URL         string  `json:"url"`
	Fields      []Field `json:"fields"`
	Color       int64   `json:"color"`
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
	return &Hook{
		Webhook: &Webhook{},
		url:     URL,
		Client:  &http.Client{Timeout: time.Second * 10},
	}
}

func (h *Hook) With(w *Webhook) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Webhook = w
}

// Run the webhook with the preconfigured settings
func (h *Hook) Run() error {
	buf := &bytes.Buffer{}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	err := json.NewEncoder(buf).Encode(h.Webhook)
	if err != nil {
		return fmt.Errorf("error encoding the webhook: %w", err)
	}

	// Building the request
	req, err := http.NewRequest(http.MethodPost, h.url, nil)
	if err != nil {
		return fmt.Errorf("error while building the request: %w", err)
	}

	if len(h.Webhook.Files) < 1 {
		req.Header.Add("content-type", "application/json")
		req.Body = io.NopCloser(buf)
	} else {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		mw.WriteField("payload_json", buf.String())

		for i, f := range h.Webhook.Files {
			ff, CFerr := mw.CreateFormFile(fmt.Sprintf("file%d", i), f.Filename)
			if CFerr != nil {
				return err
			}

			if _, CopyErr := io.Copy(ff, f.Body); CopyErr != nil {
				return CopyErr
			}
		}

		mw.Close()

		req.Header.Set("Content-Type", mw.FormDataContentType())
		req.Body = io.NopCloser(body)
	}

	// Sending the rq
	resp, err := h.Client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: error when sending: %w", err)
	}
	resp.Body.Close()

	return nil
}
