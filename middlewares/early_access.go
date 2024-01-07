package middlewares

import (
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func verifyEarlyAccessToken(tokenString string, email string) error {
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return []byte(initializers.CONFIG.EARLY_ACCESS_SECRET), nil
	// })
	// if err != nil {
	// 	return &fiber.Error{Code: 401, Message: "Invalid token. Please Generate a new one."}
	// }

	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	// if float64(time.Now().Unix()) > claims["exp"].(float64) {
	// 	return &fiber.Error{Code: 401, Message: "Your token has expired. Please Generate a new one."}
	// }

	// userEmail, ok := claims["sub"].(string)
	// if !ok {
	// 	return &fiber.Error{Code: 401, Message: "Invalid token. Please Generate a new one."}
	// }

	// if userEmail != email {
	// 	return &fiber.Error{Code: 401, Message: "Token and Email mismatch."}
	// }

	var earlyAccessModel models.EarlyAccess
	if err := initializers.DB.First(&earlyAccessModel, "email = ? AND token = ?", email, tokenString).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 401, Message: "Invalid Token."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if time.Now().After(earlyAccessModel.ExpirationTime) {
		return &fiber.Error{Code: 401, Message: "Your token has expired. Please Generate a new one."}
	}

	// if earlyAccessModel.Token != tokenString {
	// 	return &fiber.Error{Code: 401, Message: "Invalid token. Please Generate a new one."}
	// }

	return nil
	// } else {
	// 	return &fiber.Error{Code: 401, Message: "Invalid Token"}
	// }
}

func EarlyAccessCheck(c *fiber.Ctx) error {
	token := c.Query("token", "")
	if token == "" {
		return &fiber.Error{Code: 401, Message: "Provide the early access token"}
	}

	var reqBody schemas.UserCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request"}
	}

	err := verifyEarlyAccessToken(token, reqBody.Email)
	if err != nil {
		return err
	}

	return c.Next()
}
