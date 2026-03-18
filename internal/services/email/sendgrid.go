package email

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const fromEmail = "noreply@babycare.app"

// EmailService sends transactional emails via SendGrid.
type EmailService struct {
	apiKey string
}

// NewEmailService constructs an EmailService with the provided SendGrid API key.
func NewEmailService(apiKey string) *EmailService {
	return &EmailService{apiKey: apiKey}
}

// SendEmail delivers a plain-text email to a single recipient.
func (s *EmailService) SendEmail(toEmail, toName, subject, body string) error {
	from := mail.NewEmail("BabyCare", fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, body, "")

	client := sendgrid.NewSendClient(s.apiKey)
	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("email: send failed: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("email: sendgrid returned status %d: %s", resp.StatusCode, resp.Body)
	}
	return nil
}
