package validators

import (
	"strings"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
)

func EmailValidator(email string) error {
	if len(strings.Split(email, "@")) < 2 {
		return &fiber.Error{Code: 400, Message: "Invalid Email"}
	}

	domain := strings.Split(email, "@")[1]
	if !utils.Contains(config.VALID_DOMAINS, domain) {
		return &fiber.Error{Code: 400, Message: "This email is not accepted."}
	}

	return nil
}
