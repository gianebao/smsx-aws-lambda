package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/darylnwk/sms"
	"github.com/darylnwk/smsx-aws-lambda/app"
	"github.com/darylnwk/smsx-aws-lambda/app/model"
	"github.com/gianebao/shorten"
	"github.com/stretchr/testify/assert"
)

func TestRequestHandler(t *testing.T) {
	model.Shortener = shorten.Bitly{
		Username: "somebitlyuser",
		Password: `QWEUIYQWEIUASGDHA2323`,
	}

	model.Gateway = sms.Nexmo{
		APIKey:    "abcd1234",
		APISecret: "abcd1234WXYZ7890",
		From:      "Sender",
	}

	atServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "111111111111111111111111111111111111111")
	}))

	defer atServer.Close()

	shortenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "http://bit.ly/2IV3kYV")
	}))

	defer shortenServer.Close()

	shorten.BitlyAccessTokenEndpoint = atServer.URL
	shorten.BitlyEndpoint = shortenServer.URL

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

	req := events.APIGatewayProxyRequest{
		Body: `{"to": "+6599887711","text": "Hi %s! To verify, visit: %s","tokens": ["Juan",{"short-url": "https://commons.wikimedia.org/wiki/File:Smiley.svg"}]}`,
	}
	res, err := app.RequestHandler(req)

	assert.NotNil(t, req)
	assert.NotNil(t, res)
	assert.Nil(t, err)

	req = events.APIGatewayProxyRequest{
		Body: `{"to": "+6599887711","text": "Hi SMS User"}`,
	}
	res, err = app.RequestHandler(req)

	assert.NotNil(t, req)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestRequestHandler_Failed(t *testing.T) {
	model.Shortener = shorten.Bitly{
		Username: "somebitlyuser",
		Password: `QWEUIYQWEIUASGDHA2323`,
	}

	model.Gateway = sms.Nexmo{
		APIKey:    "abcd1234",
		APISecret: "abcd1234WXYZ7890",
		From:      "rdp",
	}

	atServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "111111111111111111111111111111111111111")
	}))

	defer atServer.Close()

	shortenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "http://bit.ly/2IV3kYV")
	}))

	defer shortenServer.Close()

	shorten.BitlyAccessTokenEndpoint = atServer.URL
	shorten.BitlyEndpoint = shortenServer.URL

	nexmoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message-count": "1",
          "messages": [{
              "status": "4",
              "error-text": "Bad Request"
          }]}`)
	}))

	sms.NexmoEndpoint = nexmoServer.URL

	req := events.APIGatewayProxyRequest{
		Body: `BAD JSON STRING`,
	}
	res, _ := app.RequestHandler(req)

	assert.Equal(t, `{"code":400,"text":"INVALID_REQUEST_BODY"}`, res.Body)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	req = events.APIGatewayProxyRequest{
		Body: `{"to": "+6599887711","text": "Hi %s! To verify, visit: %s","tokens": ["Juan",{"short-url": "https://commons.wikimedia.org/wiki/File:Smiley.svg"}]}`,
	}
	res, _ = app.RequestHandler(req)

	assert.Equal(t, `{"code":500,"text":"SMS_GATEWAY_ERROR[4 Bad Request]"}`, res.Body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	sms.NexmoEndpoint = "INVALIDURL"

	req = events.APIGatewayProxyRequest{
		Body: `{"to": "+6599887711","text": "Hi %s! To verify, visit: %s","tokens": ["Juan",{"short-url": "https://commons.wikimedia.org/wiki/File:Smiley.svg"}]}`,
	}
	res, _ = app.RequestHandler(req)

	assert.Equal(t, `{"code":500,"text":"SMS_GATEWAY_ERROR"}`, res.Body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
