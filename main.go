package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/darylnwk/sms"
	"github.com/darylnwk/smsx-aws-lambda/app"
	"github.com/darylnwk/smsx-aws-lambda/app/model"
	"github.com/gianebao/shorten"
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
