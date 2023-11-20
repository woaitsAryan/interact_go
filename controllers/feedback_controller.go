package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddFeedback(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	type ReqBody struct {
		Type    int    `json:"type" validate:"required"`
		Content string `json:"content" validate:"max=1000"`
	}

	var reqBody ReqBody
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[ReqBody](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	feedback := models.Feedback{
		UserID:  parsedLoggedInUserID,
		Type:    reqBody.Type,
		Content: reqBody.Content,
	}

	result := initializers.DB.Create(&feedback)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	go routines.LogFeedback(&feedback)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Feedback Submitted",
	})
}
