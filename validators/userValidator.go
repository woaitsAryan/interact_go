package validators

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UserCreateValidator(c *fiber.Ctx) error {
	var reqBody models.UserCreateSchema

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := utils.Validate[models.UserCreateSchema](reqBody); err != nil {
		return err
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.ID != uuid.Nil {
		return &fiber.Error{Code: 400, Message: "User with this Email ID already exists"}
	}

	return c.Next()
}
