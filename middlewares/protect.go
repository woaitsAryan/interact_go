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
)

func verifyToken(tokenString string, user *models.User) error { //! Crashes when both access and refresh tokens are expired
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(initializers.CONFIG.JWT_SECRET), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return &fiber.Error{Code: 403, Message: "Your token has expired."}
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			// return "", &fiber.Error{Code: 401, Message: "Invalid user ID in token claims."}
			return &fiber.Error{Code: 401, Message: "Invalid user ID in token claims."}
		}
		// return userID, nil

		initializers.DB.First(user, "id = ?", userID)

		if user.ID == uuid.Nil {
			return &fiber.Error{Code: 401, Message: "User of this token no longer exists"}
		}

		// if user.PasswordChangedAt.After(time.Unix(int64(claims["crt"].(float64)), 0)) {
		// 	return &fiber.Error{Code: 401, Message: "Password was recently changed, log in again."}
		// }

		return nil
	} else {
		return &fiber.Error{Code: 403, Message: "Invalid Token"}
	}
}

func Protect(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	tokenArr := strings.Split(authHeader, " ")

	if len(tokenArr) != 2 {
		return &fiber.Error{Code: 401, Message: "You are Not Logged In."}
	}

	tokenString := tokenArr[1]

	var user models.User
	err := verifyToken(tokenString, &user)
	if err != nil {
		return err
	}

	c.Set("loggedInUserID", user.ID.String())

	return c.Next()

}

func PartialProtect(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenArr := strings.Split(authHeader, " ")

	if len(tokenArr) != 2 {
		return c.Next()
	}

	tokenString := tokenArr[1]
	// userID, err := verifyToken(tokenString)
	// if err != nil {
	// 	return nil
	// }

	var user models.User
	err := verifyToken(tokenString, &user)
	if err != nil {
		return err
	}

	c.Set("loggedInUserID", user.ID.String())
	// c.Set("loggedInUserID", userID)

	// c.Locals("loggedInUser", user)

	return c.Next()
}
