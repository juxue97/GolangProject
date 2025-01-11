package mailer

import (
	"github.com/juxue97/common"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) (*SendGridMailer, error) {
	if apiKey == "" {
		return &SendGridMailer{}, common.ErrApiKey
	}
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}, nil
}

func (sgm *SendGridMailer) Send(templateFile, username, email string, data any, isSandBox bool) (int, error) {
	from := mail.NewEmail(fromName, sgm.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	subject, body, err := templateParsingAndBuilding(templateFile, data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandBox,
		},
	})

	err = retryWithBackoff(func() error {
		_, err := sgm.client.Send(message)
		if err != nil {
			return err
		}
		return nil
	}, maxRetries)
	if err != nil {
		return -1, err
	}

	return 200, nil
}
