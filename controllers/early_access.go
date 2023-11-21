package controllers

import (
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetEarlyAccessToken(c *fiber.Ctx) error {
	type ReqBody struct {
		Email string `json:"email" validate:"email"`
	}

	var reqBody ReqBody
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request"}
	}

	if err := helpers.Validate[ReqBody](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	var existingUser models.User
	if err := initializers.DB.First(&existingUser, "email = ?", reqBody.Email).Error; err == nil {
		return &fiber.Error{Code: 401, Message: "This email is already in use."}
	}

	// access_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": reqBody.Email,
	// 	"crt": time.Now().Unix(),
	// 	"exp": time.Now().Add(config.EARLY_ACCESS_TOKEN_TTL).Unix(),
	// })

	// access_token, err := access_token_claim.SignedString([]byte(initializers.CONFIG.EARLY_ACCESS_SECRET))
	// if err != nil {
	// 	go helpers.LogServerError("Error while decrypting Early Access Token.", err, c.Path())
	// 	return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	// }

	access_token := generateRandomString(32)

	var eaModel models.EarlyAccess

	var earlyAccessModel models.EarlyAccess
	if err := initializers.DB.First(&earlyAccessModel, "email = ?", reqBody.Email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			e := models.EarlyAccess{
				Email:          reqBody.Email,
				Token:          access_token,
				MailSent:       false,
				ExpirationTime: time.Now().Add(config.EARLY_ACCESS_TOKEN_TTL),
			}
			result := initializers.DB.Create(&e)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}

			eaModel = e
		} else {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if earlyAccessModel.CreatedAt.Add(time.Hour * 24).After(time.Now()) {
			return &fiber.Error{Code: 401, Message: "Token request limit: only once per day"}
		}
		earlyAccessModel.Token = access_token
		earlyAccessModel.ExpirationTime = time.Now().Add(config.EARLY_ACCESS_TOKEN_TTL)

		result := initializers.DB.Save(&earlyAccessModel)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
		}

		eaModel = earlyAccessModel
	}

	// err := helpers.SendMail(config.EARLY_ACCESS_EMAIL_SUBJECT, config.EARLY_ACCESS_EMAIL_BODY+access_token, "Interact User", reqBody.Email, "<div><strong>This is Valid for next 7 days!</strong></div><a href="+initializers.CONFIG.FRONTEND_URL+"/signup"+">Click Here to complete your signup!</a>")
	// if err != nil {
	// 	return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	// }

	err := helpers.SendEarlyAccessMail("Interact User", reqBody.Email, access_token)
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	eaModel.MailSent = true

	result := initializers.DB.Save(&eaModel)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Token Sent",
	})
}
