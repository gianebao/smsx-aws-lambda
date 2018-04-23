package model

import (
	"encoding/json"

	"github.com/gianebao/shorten"
	"github.com/gianebao/sms"
)

// Message represents the message payload
type Message struct {
	To       string
	Message  sms.Message
	Callback string
}

var (
	// Shortener contains the shorten.Shortener implementation from `github.com/gianebao/shorten`
	Shortener shorten.Shortener

	// Gateway contains the sms.Gateway implementation from `github.com/gianebao/sms`
	Gateway sms.Gateway
)

func processToken(m map[string]interface{}) (string, error) {
	v, ok := m["short-url"]
	if ok {
		return shorten.Shorten(Shortener, v.(string))
	}

	return v.(string), nil
}

// UnmarshalJSON parses a JSON string to a Message object
func (m *Message) UnmarshalJSON(b []byte) error {
	var (
		dat  interface{}
		err  = json.Unmarshal(b, &dat)
		datm map[string]interface{}
	)

	if err != nil {
		return err
	}

	datm = dat.(map[string]interface{})

	if _, ok := datm["to"]; ok {
		m.To = datm["to"].(string)
	}

	if _, ok := datm["callback"]; ok {
		m.Callback = datm["callback"].(string)
	}

	if _, ok := datm["text"]; ok {
		m.Message.Template = datm["text"].(string)
	}

	if v, ok := datm["tokens"]; !ok || nil == v {
		return nil
	}

	for _, v := range datm["tokens"].([]interface{}) {
		if vv, ok := v.(map[string]interface{}); ok {
			t, _ := processToken(vv)
			m.Message.Tokens = append(m.Message.Tokens, t)
		} else {
			m.Message.Tokens = append(m.Message.Tokens, v)
		}
	}

	return nil
}

// Send sends the Message as an SMS
func (m Message) Send() (sms.NexmoResponse, error) {
	var (
		smsMsg, err = sms.Send(Gateway, m.To, m.Message, m.Callback)
	)

	if err != nil {
		return sms.NexmoResponse{}, err
	}

	return smsMsg.(sms.NexmoResponse), nil
}
