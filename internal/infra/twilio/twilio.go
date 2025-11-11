package twilio

import (
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/twilio/twilio-go"
)

func NewTwilio(cfg *config.TwilioConfig) *twilio.RestClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})
	return client
}
