package types

import (
	"time"
)

type MailConfig struct {
	Exp       time.Duration
	FromEmail string
	SendGrid  SendGridConfig
	MailTrap  MailTrapConfig
}

type SendGridConfig struct {
	ApiKey string
}

type MailTrapConfig struct {
	ApiKey string
}
