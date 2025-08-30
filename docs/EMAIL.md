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

## Verify email

Create a new table verify_emails, and add a foreign key in the user table. Apply the migration to the postgres DB.

As usual, run `make sqlc` and `make mockgen`.

First a DB record is created for the verify email:

```go
verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
	Username:   user.Username,
	Email:      user.Email,
	SecretCode: util.RandomString(32),
})
if err != nil {
	return fmt.Errorf("failed to create verify email: %w", err)
}
```

Now, invoke the `SendEmail` function:

```go
subject := "Welcome to Simple Bank"
verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
	verifyEmail.ID, verifyEmail.SecretCode)
content := fmt.Sprintf(`Hello %s,<br/>Thank you for registering with us!<br/>Please <a href="%s">click here</a> to verify your email address.`, user.FullName, verifyUrl)
to := []string{user.Email}

err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
if err != nil {
	return fmt.Errorf("failed to send verify email: %w", err)
}
```

> The `mailer` attribute needs to be added to the `RedisTaskProcessor` struct.

## Create verify email RPC

Define the `rpc_verify_email.proto` file:

```proto
syntax = "proto3";

package pb;

option go_package = "simple-bank/pb";

message VerifyEmailRequest {
    int64 id = 1;
    string secret_code = 2;
}

message VerifyEmailResponse {
    bool is_verified = 1;
}
```

Now create the RPC function in the `service_simple_bank.proto`

```proto
rpc VerifyEmail(VerifyEmailRequest) returns (VerifyEmailResponse) {
	option (google.api.http) = {
	get: "/v1/verify_email"
	};
	option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	description: "This API verifies an email using gRPC"
	summary: "Verify Email"
	tags: "verify_email"
	};
}
```

Now run the `make proto` command, to generate the `Go` code.

Since we are updating the user table and verify table, its better to run these queries as a transaction. 

Please refer [this](../db/sqlc/tx_verify_email.go) file.
