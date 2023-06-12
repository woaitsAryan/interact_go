package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func Protect(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	tokenArr := strings.Split(authHeader, " ")

	if len(tokenArr) != 2 {
		return &fiber.Error{Code: 401, Message: "You are Not Logged In."}
	}

	tokenString := tokenArr[1]

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(viper.GetString("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return &fiber.Error{Code: 401, Message: "Your token has expired, log in again."}
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == uuid.Nil {
			return &fiber.Error{Code: 401, Message: "User of this token no longer exists"}
		}

		if user.PasswordChangedAt.After(time.Unix(int64(claims["exp"].(float64)), 0)) {
			return &fiber.Error{Code: 401, Message: "Password was recently changed, log in again."}
		}

		c.Locals("user", user)

		return c.Next()

	} else {
		return &fiber.Error{Code: 401, Message: "Invalid Token"}
	}

}
