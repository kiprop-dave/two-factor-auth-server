package services

import (
	"log"

	config "github.com/kiprop-dave/2faAuth/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(target, body string) error {
	env := config.Environment
	sid := env.TwilioSid
	token := env.TwilioToken
	phone := env.PhoneNumber

	twilioParams := twilio.ClientParams{
		Username: sid,
		Password: token,
	}
	client := twilio.NewRestClientWithParams(twilioParams)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(target)
	params.SetFrom(phone)
	params.SetBody(body)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
