package controllers

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func SharePost(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content string         `json:"content"`
		Chats   pq.StringArray `json:"chats"`
		PostID  string         `json:"postID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chats := reqBody.Chats

	for _, chatID := range chats {
		message := models.Message{
			UserID:  parsedUserID,
			Content: reqBody.Content,
		}

		if reqBody.PostID != "" {
			parsedPostID, err := uuid.Parse(reqBody.PostID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
			}
			message.PostID = &parsedPostID

			parsedChatID, err := uuid.Parse(chatID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID."}
			}

			message.ChatID = parsedChatID

			result := initializers.DB.Create(&message)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			go routines.IncrementPostShare(parsedPostID)
		} else {
			return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
		}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Shared",
	})

}

func ShareProject(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content   string         `json:"content"`
		Chats     pq.StringArray `json:"chats"`
		ProjectID string         `json:"projectID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chats := reqBody.Chats

	for _, chatID := range chats {
		message := models.Message{
			UserID:  parsedUserID,
			Content: reqBody.Content,
		}

		if reqBody.ProjectID != "" {
			parsedProjectID, err := uuid.Parse(reqBody.ProjectID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
			}
			message.ProjectID = &parsedProjectID

			parsedChatID, err := uuid.Parse(chatID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID."}
			}

			message.ChatID = parsedChatID

			result := initializers.DB.Create(&message)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
			}

			go routines.IncrementProjectShare(parsedProjectID)

		} else {
			return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
		}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project Shared",
	})

}
