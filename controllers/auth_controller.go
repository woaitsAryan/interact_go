package controllers

import (
	"fmt"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createSendToken(c *fiber.Ctx, user models.User, statusCode int, message string) error {
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
	})

	return c.Status(statusCode).JSON(fiber.Map{
		"status":     "success",
		"message":    message,
		"token":      access_token,
		"userID":     user.ID,
		"profilePic": user.ProfilePic,
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

	newUser := models.User{
		Name:              reqBody.Name,
		Email:             reqBody.Email,
		Password:          string(hash),
		Username:          reqBody.Username,
		PasswordChangedAt: time.Now(),
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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

func Refresh(c *fiber.Ctx) error {

	var reqBody struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	access_token_string := reqBody.Token

	access_token, _ := jwt.Parse(access_token_string, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(initializers.CONFIG.JWT_SECRET), nil
	})

	if access_token_claims, ok := access_token.Claims.(jwt.MapClaims); ok {

		access_token_userID, ok := access_token_claims["sub"].(string)
		if !ok {
			return &fiber.Error{Code: 401, Message: "Invalid user ID in token claims."}
		}

		var user models.User
		err := initializers.DB.First(&user, "id = ?", access_token_userID).Error
		if err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		if user.ID == uuid.Nil {
			return &fiber.Error{Code: 401, Message: "User of this token no longer exists"}
		}

		refresh_token_string := c.Cookies("refresh_token")
		if refresh_token_string == "" {
			return &fiber.Error{Code: 401, Message: "Session Expired, Log In Again"}
		}

		refresh_token, _ := jwt.Parse(refresh_token_string, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(initializers.CONFIG.JWT_SECRET), nil
		})

		if refresh_token_claims, ok := refresh_token.Claims.(jwt.MapClaims); ok && refresh_token.Valid {
			refresh_token_userID, ok := refresh_token_claims["sub"].(string)
			if !ok {
				return &fiber.Error{Code: 401, Message: "Invalid user ID in token claims."}
			}

			if refresh_token_userID != access_token_userID {
				return &fiber.Error{Code: 401, Message: "Mismatched Tokens."}
			}

			if time.Now().After(time.Unix(int64(refresh_token_claims["exp"].(float64)), 0)) {
				return &fiber.Error{Code: 401, Message: "Token has expired, Log in again."}
			}

			new_access_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.ID,
				"crt": time.Now().Unix(),
				"exp": time.Now().Add(config.ACCESS_TOKEN_TTL).Unix(),
			})

			new_access_token, err := new_access_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
			if err != nil {
				go helpers.LogServerError("Error while decrypting JWT Token.", err, c.Path())
				return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
			}

			return c.Status(200).JSON(fiber.Map{
				"status": "success",
				"token":  new_access_token,
			})
		}

		return nil

	} else {
		return &fiber.Error{Code: 401, Message: "Invalid Token"}
	}
}
