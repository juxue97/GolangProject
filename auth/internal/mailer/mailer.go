package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
	"time"

	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/common"
)

const (
	fromName                      = "Mr. Teh"
	maxRetries                    = 3
	UserWelcomeInvitationTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var fs embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandBox bool) (int, error)
}

// var SendGridMailer *SendGridMailer
var (
	MailTrapMailer MailTrapClient
	err            error
)

func init() {
	// SendGridMailer,err = NewSendGrid(config.Configs.Mail.SendGrid.ApiKey, config.Configs.Mail.FromEmail)
	MailTrapMailer, err = NewMailTrap(config.Configs.Mail.MailTrap.ApiKey, config.Configs.Mail.FromEmail)
	if err != nil {
		common.Logger.Fatal(err)
	}
	common.Logger.Info("Mailer Service initialized")
}

func templateParsingAndBuilding(templateFile string, data any) (*bytes.Buffer, *bytes.Buffer, error) {
	// template parsing and building
	tmpl, err := template.ParseFS(fs, "templates/"+templateFile)
	if err != nil {
		return nil, nil, err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, nil, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil, nil, err
	}

	return subject, body, nil
}

func retryWithBackoff(operation func() error, maxRetries int) error {
	var retryErr error
	for i := 0; i < maxRetries; i++ {
		retryErr = operation()
		if retryErr == nil {
			return nil // Success
		}

		// Optional backoff
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("operation failed after %d attempts, last error: %v", maxRetries, retryErr)
}
