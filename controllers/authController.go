package controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createSendToken(c *fiber.Ctx, user models.User, statusCode int, message string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	tokenString, err := token.SignedString([]byte(initializers.CONFIG.JWT_SECRET))

	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}

	//set cookie

	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"token":   tokenString,
		"userID":  user.ID,
	})
}

func SignUp(c *fiber.Ctx) error {

	var reqBody schemas.UserCreateSchema

	c.BodyParser(&reqBody)

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)

	if err != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error"}
	}

	newUser := models.User{
		Name:              reqBody.Name,
		Email:             reqBody.Email,
		Password:          string(hash),
		Username:          reqBody.Username,
		PasswordChangedAt: time.Now(),
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating user"}
	}

	return createSendToken(c, newUser, 201, "Account Created")
}

func LogIn(c *fiber.Ctx) error {

	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var user models.User

	initializers.DB.First(&user, "username = ?", reqBody.Username)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user with these credentials found."}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "No user with these credentials found."}
	}

	return createSendToken(c, user, 200, "Logged In")
}
