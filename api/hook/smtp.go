package hook

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/octacian/backroom/api/cage"
	"github.com/octacian/backroom/api/config"
	"github.com/wneessen/go-mail"
)

// SMTPAdapter is an adapter that sends an email.
type SMTPAdapter struct {
	client *mail.Client
}

// NewSMTPAdapter creates a new SMTPAdapter with the configured SMTP client.
// See config.RC.Mail for configuration.
func NewSMTPAdapter() (*SMTPAdapter, error) {
	opts := make([]mail.Option, 0)
	if config.RC.Mail.SMTP.TLS {
		opts = append(opts, mail.WithSSL())
	}
	if config.RC.Mail.SMTP.Port != 0 {
		opts = append(opts, mail.WithPort(config.RC.Mail.SMTP.Port))
	}
	opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
	opts = append(opts, mail.WithUsername(config.RC.Mail.SMTP.Username))
	opts = append(opts, mail.WithPassword(config.RC.Mail.SMTP.Password))

	client, err := mail.NewClient(config.RC.Mail.SMTP.Host, opts...)
	if err != nil {
		return nil, err
	}

	return &SMTPAdapter{client: client}, nil
}

// FromAddress returns the configured from address for the SMTP client.
func (a *SMTPAdapter) FromAddress() string {
	fromName := config.RC.Mail.FromName
	if fromName == "" {
		fromName = config.RC.Mail.FromAddress
	}
	if fromName == "" {
		fromName = "backroom"
	}
	return fmt.Sprintf("%s <%s>", fromName, config.RC.Mail.FromAddress)
}

// Run executes the SMTPAdapter with the given hook and record.
func (a *SMTPAdapter) Run(hook *Hook, record *cage.Record) error {
	message := mail.NewMsg()

	if err := message.From(a.FromAddress()); err != nil {
		return err
	}

	if err := message.To(hook.Target); err != nil {
		return err
	}

	jsonText, err := json.MarshalIndent(record.Data, "", "  ")
	if err != nil {
		return err
	}

	message.Subject(fmt.Sprintf("%s: new record created", record.Cage))
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf("%s record %s created\n\n%s", record.Cage, record.UUID, jsonText))

	if err := a.client.DialAndSend(message); err != nil {
		return err
	}

	slog.Info("SMTPAdapter sent email", "to", hook.Target, "uuid", record.UUID)
	return nil
}
