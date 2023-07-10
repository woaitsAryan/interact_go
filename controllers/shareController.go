package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SharePost(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content string `json:"content"`
		ChatID  string `json:"chatID"`
		PostID  string `json:"postID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

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

		var post models.Post
		if err := initializers.DB.First(&post, "id=?", parsedPostID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}

		post.NoShares++

		result := initializers.DB.Save(post)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while updating Post."}
		}

		parsedChatID, err := uuid.Parse(chatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		message.ChatID = parsedChatID

		result = initializers.DB.Create(&message)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the message."}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Post Shared",
		})

	} else {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
	}
}

func ShareProject(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content   string `json:"content"`
		ChatID    string `json:"chatID"`
		ProjectID string `json:"projectID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

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

		var project models.Project
		if err := initializers.DB.First(&project, "id=?", parsedProjectID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}

		project.NoShares++

		result := initializers.DB.Save(project)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Database Error while updating Project."}
		}

		parsedChatID, err := uuid.Parse(chatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		message.ChatID = parsedChatID

		result = initializers.DB.Create(&message)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the message."}
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Project Shared",
		})

	} else {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
	}
}
