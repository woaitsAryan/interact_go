package helpers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/go-gomail/gomail"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/sendgrid/sendgrid-go"
	// "github.com/sendgrid/sendgrid-go/helpers/mail"
)

//* SendGrid Mailer
// func SendMail(subject string, body string, recipientName string, recipientEmail string, htmlStr string) error {
// 	from := mail.NewEmail(config.SENDER_NAME, config.SENDER_EMAIL)
// 	to := mail.NewEmail(recipientName, recipientEmail)
// 	htmlContent := body + htmlStr
// 	message := mail.NewSingleEmail(from, subject, to, body, htmlContent)
// 	client := sendgrid.NewSendClient(initializers.CONFIG.SENDGRID_KEY)
// 	_, err := client.Send(message)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func SendMail(subject string, body string, recipientName string, recipientEmail string, htmlStr string) error {
	htmlContent := body + htmlStr //TODO22 Email Template
	m := gomail.NewMessage()
	m.SetHeader("From", config.GMAIL_SENDER)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.GMAIL_SENDER, initializers.CONFIG.GMAIL_KEY)

	if err := d.DialAndSend(m); err != nil {
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
		return
	}

	t.Execute(&body, struct {
		Name         string
		ChatUserName string
	}{Name: recipientName, ChatUserName: chatUserName})

	if err != nil {
		LogDatabaseError("Error while sending Chat Mail-SendChatMail", err, "go_routine")
		return
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

func createMailerJWT() (string, error) {
	token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "backend",
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.LOGIN_TOKEN_TTL).Unix(),
	})

	token, err := token_claim.SignedString([]byte(initializers.CONFIG.MAILER_SECRET))
	if err != nil {
		return "", err
	}

	return token, nil
}

func SendMailReq(email string, emailType int, user *models.User) error {
	jsonData, err := json.Marshal(map[string]any{
		"email": email,
		"type":  emailType,
		"user":  user,
	})
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}

	jwt, err := createMailerJWT()
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}

	request, err := http.NewRequest("POST", initializers.CONFIG.MAILER_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.Header.Set("API-TOKEN", initializers.CONFIG.MAILER_TOKEN)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}
	defer response.Body.Close()

	return nil
}
