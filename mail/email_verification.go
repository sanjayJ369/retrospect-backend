package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/sanjayj369/retrospect-backend/token"
)

func SendVerificationMail(
	sender EmailSender,
	userId uuid.UUID,
	to string,
	tokenMaker token.Maker,
	duration time.Duration,
	endpoint string,
	templateFile string) error {
	tkn, _, err := tokenMaker.CreateToken(userId, duration)
	if err != nil {
		return fmt.Errorf("unable to create verification token: %w", err)
	}

	verficationLink := fmt.Sprintf("%s?token=%s", endpoint, tkn)
	tmp, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("parsing email template failed: %w", err)
	}

	content := bytes.NewBufferString("")
	err = tmp.Execute(content, map[string]string{
		"VerificationLink": verficationLink,
	})
	if err != nil {
		return fmt.Errorf("executing email template failed: %w", err)
	}

	subject := "Verify your email address"
	err = sender.SendMail(subject, content.String(), []string{to}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("sending email failed: %w", err)
	}

	return nil
}
