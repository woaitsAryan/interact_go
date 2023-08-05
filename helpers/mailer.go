package helpers

import (
	"context"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/mailgun/mailgun-go/v4"
)

func SendMail(subject string, body string, recipient string) error {
	sender := config.SENDER + "@" + config.MAILGUN_DOMAIN

	mg := mailgun.NewMailgun(config.MAILGUN_DOMAIN, config.MAILGUN_API_KEY)
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
