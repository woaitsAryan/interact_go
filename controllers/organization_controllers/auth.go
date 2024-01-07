package organization_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func createSendToken(c *fiber.Ctx, user models.User, statusCode int, message string, org models.Organization) error {
	access_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.ACCESS_TOKEN_TTL).Unix(),
	})

	refresh_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.REFRESH_TOKEN_TTL).Unix(),
	})

	access_token, err := access_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
	if err != nil {
		go helpers.LogServerError("Error while decrypting JWT Token.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	refresh_token, err := refresh_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
	if err != nil {
		go helpers.LogServerError("Error while decrypting JWT Token.", err, c.Path())
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Expires:  time.Now().Add(config.REFRESH_TOKEN_TTL),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(statusCode).JSON(fiber.Map{
		"status":       "success",
		"message":      message,
		"token":        access_token,
		"user":         user,
		"email":        user.Email,
		"phoneNo":      user.PhoneNo,
		"organization": org,
	})
}

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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	organization := models.Organization{
		UserID:            newOrg.ID,
		OrganizationTitle: newOrg.Name,
		CreatedAt:         time.Now(),
	}

	result = initializers.DB.Create(&organization)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	newProfile := models.Profile{
		UserID: newOrg.ID,
	}

	result = initializers.DB.Create(&newProfile)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	go routines.SendWelcomeNotification(newOrg.ID)
	go routines.MarkOrganizationHistory(organization.ID, newOrg.ID, -1, nil, nil, nil, nil, nil, "")

	return createSendToken(c, newOrg, 201, "Organization Created", organization)
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
	if err := initializers.DB.First(&user, "email = ? AND organization_status = true", reqBody.Email).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No account with these credentials found."}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "No account with these credentials found."}
	}

	user.LastLoggedIn = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "user_id=?", user.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return createSendToken(c, user, 200, "Logged In", organization)
}

func OAuthLogIn(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var user models.User
	if err := initializers.DB.Session(&gorm.Session{SkipHooks: true}).First(&user, "id = ? AND organization_status = true", loggedInUserID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No organization with these credentials found."}
	}

	if !user.Active {
		if time.Now().After(user.DeactivatedAt.Add(30 * 24 * time.Hour)) {
			return &fiber.Error{Code: 400, Message: "Cannot Log into a deactivated account."}
		}
		user.Active = true
	}

	user.LastLoggedIn = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var organization models.Organization
	if err := initializers.DB.First(&organization, "user_id=?", user.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return createSendToken(c, user, 200, "Logged In", organization)
}
