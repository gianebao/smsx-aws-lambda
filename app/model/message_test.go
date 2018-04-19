package model_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylnwk/sms"
	"github.com/darylnwk/smsx-aws-lambda/app/model"
	"github.com/gianebao/shorten"
	"github.com/stretchr/testify/assert"
)

func TestMessage_UnmarshalJSON(t *testing.T) {
	var m model.Message

	model.Shortener = shorten.Bitly{
		Username: "somebitlyuser",
		Password: `QWEUIYQWEIUASGDHA2323`,
	}

	atServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "111111111111111111111111111111111111111")
	}))

	defer atServer.Close()

	shortenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "http://bit.ly/a234567")
	}))

	defer shortenServer.Close()

	shorten.BitlyAccessTokenEndpoint = atServer.URL
	shorten.BitlyEndpoint = shortenServer.URL

	err := json.Unmarshal([]byte(`{
    "to": "+6591110000",
    "text": "Hi %s! To verify, visit: %s",
    "tokens": [
          "Juan",
          {"short-url": "https://commons.wikimedia.org/wiki/File:Smiley.svg"}
    ]}`),
		&m)

	assert.Equal(t, "+6591110000", m.To)
	assert.Equal(t, "Hi Juan! To verify, visit: http://bit.ly/a234567", m.Message.String())
	assert.Nil(t, err)

	m = model.Message{}
	err = json.Unmarshal([]byte(`{
		"to": "+6591110000",
		"text": "Hello world"
		}`),
		&m)

	assert.Equal(t, "+6591110000", m.To)
	assert.Equal(t, "Hello world", m.Message.String())
	assert.Nil(t, err)

	m = model.Message{}
	err = json.Unmarshal([]byte(`{INVALID JSON}`),
		&m)

	assert.EqualError(t, err, `invalid character 'I' looking for beginning of object key string`)
}

func TestMessage_Send(t *testing.T) {
	model.Gateway = sms.Nexmo{
		APIKey:    "abcd1234",
		APISecret: "abcd1234WXYZ7890",
		From:      "Sender",
	}

	nexmoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message-count": "1","messages": [{
							"to": "99887711",
							"message-id": "0F0000008BD3AD66",
							"status": "0",
							"remaining-balance": "1.97600000",
							"message-price": "0.02400000",
							"network": "52501"
					}]}`)
	}))

	sms.NexmoEndpoint = nexmoServer.URL

	m := model.Message{
		To: "99887711",
		Message: sms.Message{
			Template: "Hello World",
		},
	}

	res, err := m.Send()

	assert.Equal(t, "1", res.MessageCount)
	assert.Equal(t, "99887711", res.Messages[0].To)
	assert.Equal(t, "0F0000008BD3AD66", res.Messages[0].MessageID)
	assert.Equal(t, "0", res.Messages[0].Status)
	assert.Equal(t, "1.97600000", res.Messages[0].RemainingBalance)
	assert.Equal(t, "0.02400000", res.Messages[0].MessagePrice)
	assert.Equal(t, "52501", res.Messages[0].Network)
	assert.Nil(t, err)
}

func TestMessage_SendFailed(t *testing.T) {
	model.Gateway = sms.Nexmo{
		APIKey:    "abcd1234",
		APISecret: "abcd1234WXYZ7890",
		From:      "Sender",
	}

	nexmoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `-}`)
	}))

	sms.NexmoEndpoint = nexmoServer.URL

	m := model.Message{
		To: "99887711",
		Message: sms.Message{
			Template: "Hello World",
		},
	}

	res, err := m.Send()

	assert.Equal(t, sms.NexmoResponse{}, res)
	assert.EqualError(t, err, "invalid character '}' in numeric literal")
}
