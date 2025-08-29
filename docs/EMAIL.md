## Sending Email

Install the following plugin:

```sh
github.com/jordan-wright/email
```

Add the necessary ENV variables, in this case:

```env
EMAIL_SENDER_NAME=<name>
EMAIL_SENDER_ADDRESS=<email>
EMAIL_SENDER_PASSWORD=<app-password>
```

Now create an interface `EmailSender`. We could have created just one method for sending emails; we created an interface instead so that we can add other email clients if needed.

```go
type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string // Use an app password
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{name: name, fromEmailAddress: fromEmailAddress, fromEmailPassword: fromEmailPassword}
}
```

Now import the package `email` package:

```go
import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

func (sender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	email.Subject = subject
	email.HTML = []byte(content)
	email.To = to
	email.Cc = cc
	email.Bcc = bcc

	for _, file := range attachFiles {
		_, err := email.AttachFile(file)
		if err != nil {
			return err
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)

	return email.Send(smtpServerAddress, smtpAuth)
}
```

We can write a simple test case validate this function.

```go
package mail

import (
	"simple-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>This is a test email</p>
	`
	to := []string{"ombalapure7@gmail.com", "odbalapure@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, []string{}, []string{}, attachFiles)
	require.NoError(t, err)
}
```

> There are other email clients that require domain verification step eg: Amazon SES.

Everytime this test run, an email will be sent out, so we skip this test using the `-short` flag. Add this flag while running the test command:

```sh
go test -v -short -cover ./...
```

Also, update the test function

```go
func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
    // ...
}
```

If we take a look at the console, the test will be skipped.

```
=== RUN   TestSendEmailWithGmail
    sender_test.go:12: 
--- SKIP: TestSendEmailWithGmail (0.00s)
```
