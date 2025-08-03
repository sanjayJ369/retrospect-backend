package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMailgunSender(t *testing.T) {
	sender, err := NewMailgunSender()
	require.NoError(t, err)

	subject := "hi this is a test :)"
	content := `
		<h1>HTML ohh...</h1>
		this is my mail btw
	`
	to := []string{"j.sanjay336699@gmail.com"}
	attachments := []string{"../README.md"}

	err = sender.SendMail(subject, content, to, nil, nil, attachments)
	require.NoError(t, err)
}
