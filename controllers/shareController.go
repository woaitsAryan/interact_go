package controllers

import (
	"log"

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

		parsedChatID, err := uuid.Parse(chatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		message.ChatID = parsedChatID

		result := initializers.DB.Create(&message)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the message."}
		}

		go incrementPostShare(parsedPostID)

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

		parsedChatID, err := uuid.Parse(chatID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid ID."}
		}

		message.ChatID = parsedChatID

		result := initializers.DB.Create(&message)
		if result.Error != nil {
			return &fiber.Error{Code: 500, Message: "Internal Server Error while creating the message."}
		}

		go incrementProjectShare(parsedProjectID)

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Project Shared",
		})

	} else {
		return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
	}
}

func incrementPostShare(postID uuid.UUID) {
	var post models.Post
	if err := initializers.DB.First(&post, "id=?", postID).Error; err != nil {
		log.Println("No Post of this ID found.")
	} else {
		post.NoShares++
		result := initializers.DB.Save(post)
		if result.Error != nil {
			log.Println("Database Error while updating Post.")
		}
	}
}

func incrementProjectShare(projectID uuid.UUID) {
	var project models.Project
	if err := initializers.DB.First(&project, "id=?", projectID).Error; err != nil {
		log.Println("No Project of this ID found.")
	} else {
		project.NoShares++
		result := initializers.DB.Save(project)
		if result.Error != nil {
			log.Println("Database Error while updating Project.")
		}
	}
}
