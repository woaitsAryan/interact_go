package helpers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(subject string, body string, recipientName string, recipientEmail string, htmlStr string) error {
	from := mail.NewEmail(config.SENDER_NAME, config.SENDER_EMAIL)
	to := mail.NewEmail(recipientName, recipientEmail)
	htmlContent := body + htmlStr //TODO Email Template
	message := mail.NewSingleEmail(from, subject, to, body, htmlContent)
	client := sendgrid.NewSendClient(initializers.CONFIG.SENDGRID_KEY)
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
