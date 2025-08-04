package mail

// func TestMailgunSender(t *testing.T) {
// 	config, err := util.LoadConfig("..")
// 	require.NoError(t, err)
// 	fmt.Println(config)
// 	sender, err := NewMailgunSender(config)
// 	require.NoError(t, err)

// 	subject := "hi this is a test :)"
// 	content := `
// 		<h1>HTML ohh...</h1>
// 		this is my mail btw
// 	`
// 	to := []string{"j.sanjay336699@gmail.com"}
// 	attachments := []string{"../README.md"}

// 	err = sender.SendMail(subject, content, to, nil, nil, attachments)
// 	require.NoError(t, err)
// }
