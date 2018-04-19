package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/darylnwk/sms"
	"github.com/darylnwk/smsx-aws-lambda/app/model"
)

// RequestHandler handles the API gateway request
func RequestHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		smsRes sms.NexmoResponse
		m      model.Message
		err    error
		resp   = events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)

	if err = json.Unmarshal([]byte(req.Body), &m); err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = model.Status{
			Code: resp.StatusCode,
			Text: "INVALID_REQUEST_BODY",
		}.String()

		return resp, err
	}

	if smsRes, err = m.Send(); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = model.Status{
			Code: resp.StatusCode,
			Text: "SMS_GATEWAY_ERROR",
		}.String()

		return resp, err
	}

	if 0 < len(smsRes.Messages) && "0" == smsRes.Messages[0].Status {
		resp.Body = fmt.Sprintf(
			`{"id":"%s", "status":%s, "network":"%s"}`,
			smsRes.Messages[0].MessageID,
			smsRes.Messages[0].Status,
			smsRes.Messages[0].Network,
		)
	} else {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = model.Status{
			Code: resp.StatusCode,
			Text: fmt.Sprintf("SMS_GATEWAY_ERROR[%s %s]", smsRes.Messages[0].Status, smsRes.Messages[0].ErrorText),
		}.String()
	}

	return resp, nil
}
