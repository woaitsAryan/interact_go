package controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func createSendToken(c *fiber.Ctx, user models.User, statusCode int, message string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))

	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}

	//set cookie

	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"user":    user,
		"token":   tokenString,
	})
}

func SignUp(c *fiber.Ctx) error {

	var reqBody models.UserCreateSchema

	c.BodyParser(&reqBody)

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)

	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}

	var picName string = ""

	picName, err = utils.SaveFile(c, "profilePic", "users/profilePics", true, true)
	if err != nil {
		return err
	}

	newUser := models.User{
		Name:       reqBody.Name,
		Email:      reqBody.Email,
		Password:   string(hash),
		Username:   reqBody.Username,
		PhoneNo:    reqBody.PhoneNo,
		ProfilePic: picName,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating user"}
	}

	return createSendToken(c, newUser, 201, "Account Created")
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

	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user with these credentials found."}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "No user with these credentials found."}
	}

	return createSendToken(c, user, 200, "Logged In")
}
