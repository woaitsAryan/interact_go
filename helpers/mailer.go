package helpers

import (
	"bytes"
	"html/template"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/go-gomail/gomail"
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

func SendEarlyAccessMail(recipientName string, recipientEmail string, token string) error {
	var body bytes.Buffer
	path := config.TEMPLATE_DIR + "early_access.html"
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	t.Execute(&body, struct {
		Name  string
		Email string
		Token string
	}{Name: recipientName, Email: recipientEmail, Token: token})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.GMAIL_SENDER)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Interact Early Access")
	m.SetBody("text/html", body.String())
	// m.Attach("attachment.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, config.GMAIL_SENDER, initializers.CONFIG.GMAIL_KEY)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendChatMail(recipientName string, recipientEmail string, chatUserName string) {
	var body bytes.Buffer
	path := config.TEMPLATE_DIR + "chat.html"
	t, err := template.ParseFiles(path)
	if err != nil {
		LogDatabaseError("Error while sending Chat Mail-SendChatMail", err, "go_routine")
	}

	t.Execute(&body, struct {
		Name         string
		ChatUserName string
	}{Name: recipientName, ChatUserName: chatUserName})

	if err != nil {
		LogDatabaseError("Error while sending Chat Mail-SendChatMail", err, "go_routine")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.GMAIL_SENDER)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "You have a new chat!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, config.GMAIL_SENDER, initializers.CONFIG.GMAIL_KEY)

	if err := d.DialAndSend(m); err != nil {
		LogDatabaseError("Error while sending Chat Mail-SendChatMail", err, "go_routine")
	}
}
