package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type GroupMeMessage struct {
	BotId string `json:"bot_id"`
	Text  string `json:"text"`
}

type GroupMeCallback struct {
	Name       string `json:"name"`
	SenderType string `json:"sender_type"`
	Text       string `json:"text"`
}

type GroupMe interface {
	Post(string) (int, error)
}

type GroupMeMessenger struct {
	botId   string
	postURL *url.URL
}

func NewGroupMeMessenger(botId string, postURL *url.URL) GroupMe {
	return &GroupMeMessenger{
		botId:   botId,
		postURL: postURL,
	}
}

func (g *GroupMeMessenger) Post(text string) (int, error) {
	msg, err := json.Marshal(&GroupMeMessage{
		BotId: g.botId,
		Text:  text,
	})
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(
		g.postURL.String(), "application/json", bytes.NewBuffer(msg))
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func GroupMeParse(text string) (*GroupMeCallback, error) {
	var cb *GroupMeCallback
	err := json.Unmarshal([]byte(text), &cb)
	if err != nil {
		return nil, err
	}
	return cb, err
}
