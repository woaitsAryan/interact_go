package organization_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *fiber.Ctx) error {
	var reqBody schemas.UserCreateSchema

	c.BodyParser(&reqBody)

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		go helpers.LogServerError("Error while hashing Password.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	newOrg := models.User{
		Name:               reqBody.Name,
		Email:              reqBody.Email,
		Password:           string(hash),
		Username:           reqBody.Username,
		PasswordChangedAt:  time.Now(),
		OrganizationStatus: true,
	}

	result := initializers.DB.Create(&newOrg)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	organization := models.Organization{
		UserID:            newOrg.ID,
		OrganizationTitle: newOrg.Name,
		CreatedAt:         time.Now(),
	}

	result = initializers.DB.Create(&organization)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	go routines.SendWelcomeNotification(newOrg.ID)

	return controllers.CreateSendToken(c, newOrg, 201, "Organization Created")
}

func LogIn(c *fiber.Ctx) error {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var user models.User
	initializers.DB.Session(&gorm.Session{SkipHooks: true}).First(&user, "email = ? AND organization_status = true", reqBody.Email)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No account with these credentials found."}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "No account with these credentials found."}
	}

	user.LastLoggedIn = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return controllers.CreateSendToken(c, user, 200, "Logged In")
}
