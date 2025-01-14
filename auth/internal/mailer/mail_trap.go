package mailer

import (
	"errors"

	gomail "gopkg.in/mail.v2"
)

type MailTrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (MailTrapClient, error) {
	if apiKey == "" {
		return MailTrapClient{}, errors.New("api key is required")
	}

	return MailTrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (mtc MailTrapClient) Send(templateFile, username, email string, data any, isSandBox bool) (int, error) {
	// Template parsing and building
	subject, body, err := templateParsingAndBuilding(templateFile, data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", mtc.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", mtc.apiKey)

	err = retryWithBackoff(func() error {
		return dialer.DialAndSend(message)
	}, maxRetries)
	if err != nil {
		return -1, err
	}

	return 200, nil
}
