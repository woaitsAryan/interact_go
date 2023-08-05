package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/mailgun/mailgun-go/v4"
)

func SendMail(subject string, body string, recipient string) error {
	sender := config.SENDER + "@" + config.MAILGUN_DOMAIN

	mg := mailgun.NewMailgun(initializers.CONFIG.MAILGUN_DOMAIN, initializers.CONFIG.MAILGUN_PRIVATE_API_KEY)
	message := mg.NewMessage(sender, subject, body, recipient)

	// mg.SetAPIBase(mailgun.APIBaseEU)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mes, _, err := mg.Send(ctx, message)
	if err != nil {
		fmt.Println(mes)
		return err
	}

	return nil
}
