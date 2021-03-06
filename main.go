package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gianebao/shorten"
	"github.com/gianebao/sms"
	"github.com/gianebao/smsx-aws-lambda/app"
	"github.com/gianebao/smsx-aws-lambda/app/model"
)

var (
	bitly shorten.Bitly
	nexmo sms.Nexmo
)

// Init initializes parameters in starting the application. INIT has to be manually executed for lambda
func Init() {
	nexmo.APIKey = os.Getenv("NEXMOAPIKEY")
	nexmo.APISecret = os.Getenv("NEXMOAPISECRET")
	nexmo.From = os.Getenv("NEXMOFROM")
	bitly.Username = os.Getenv("BITLYUSERNAME")
	bitly.Password = os.Getenv("BITLYPASSWORD")
}

func main() {
	Init()

	model.Shortener = bitly
	model.Gateway = nexmo

	lambda.Start(app.RequestHandler)
}
