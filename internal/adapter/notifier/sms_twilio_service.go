package notifier

import (
	"context"
	"log/slog"
	"strings"

	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSNotifier struct {
	Config *config.TwilioConfig
	Client *twilio.RestClient
}

func NewSMSNotifier(twilioClient *twilio.RestClient, config *config.TwilioConfig) *SMSNotifier {
	return &SMSNotifier{
		Client: twilioClient,
		Config: config,
	}
}

// NotifySMS sends an SMS via Twilio
func (s *SMSNotifier) NotifySMS(ctx context.Context, phone string, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetMessagingServiceSid(s.Config.MessageServiceSID)
	params.SetBody(message)

	_, err := s.Client.Api.CreateMessage(params)
	if err != nil {
		// Check for authentication errors (20003)
		errMsg := err.Error()
		if containsAny(errMsg, "20003", "Authenticate") {
			slog.Error("Twilio authentication failed - Check your credentials",
				slog.String("error", errMsg),
				slog.String("account_sid", s.Config.AccountSID[:10]+"..."),
				slog.String("help", "Verify TWILIO_ACCOUNT_SID starts with 'AC' at https://console.twilio.com/"))
		} else {
			slog.Error("Error sending SMS via Twilio", "error", err)
		}
		return err
	}

	return nil
}

// containsAny checks if the string contains any of the substrings
func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
// NotifySMSMulti sends SMS messages to multiple recipients via Twilio
func (s *SMSNotifier) NotifySMSMulti(ctx context.Context, phones []string, message string) error {
	// perfomance can be improved by using goroutines for concurrent sending
	for i := 0; i < len(phones); i += MaxBatchSize {
		end := i + MaxBatchSize
		if end > len(phones) {
			end = len(phones)
		}
		batch := phones[i:end]
		go func(batch []string) {
			for _, phone := range batch {
				err := s.NotifySMS(ctx, phone, message)
				if err != nil {
					slog.Error("Error sending SMS to "+phone, "error", err)
					// continue sending to other numbers even if one fails
				}
			}
		}(batch)
	}
	return nil
}
