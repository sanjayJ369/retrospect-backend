package mail

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/sanjayj369/retrospect-backend/util"
)

type EmailSender interface {
	SendMail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type MailgunSender struct {
	mg     mailgun.Mailgun
	domain string
	sender string
}

func NewMailgunSender(config util.Config) (EmailSender, error) {
	mg := mailgun.NewMailgun(config.MailgunAPIKEY)

	return &MailgunSender{
		mg:     mg,
		domain: config.MailgunDomain,
		sender: fmt.Sprintf(`"Retrospect" <noreplay@%s>`, config.MailgunDomain),
	}, nil
}

func (sender *MailgunSender) SendMail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	message := mailgun.NewMessage(sender.domain, sender.sender, subject, "", to...)

	for _, i := range cc {
		message.AddCC(i)
	}
	for _, i := range bcc {
		message.AddBCC(i)
	}

	message.SetHTML(content)

	for _, file := range attachFiles {
		message.AddAttachment(file)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	resp, err := sender.mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("Email sent successfully!, Response: %s\n", resp)
	return nil
}
