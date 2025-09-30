package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apikey    string
	client    *sendgrid.Client
}

func NewSendGrid(fromEmail, apikey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apikey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apikey:    apikey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < MaxRetries; i++ {
		res, err := s.client.Send(message)
		if err != nil {
			log.Printf("faliled to send mail attempt %d out of %d", i+1, MaxRetries)
			continue
		}
		log.Printf("successfuly send mail ot %v with status code %v", email, res.StatusCode)
		log.Print("result", res)
		return nil
	}
	fmt.Errorf("failed to send email after %d attempts", MaxRetries)
	return nil
}
