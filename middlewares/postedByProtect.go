package middlewares

import (
	"errors"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PostUserProtect(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	postID := c.Params("postID")

	var post models.Post

	if err := initializers.DB.First(&post, "id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	if post.UserID.String() != userID {
		return &fiber.Error{Code: 403, Message: "Not Allowed to Perfom this Action."}
	}

	return c.Next()
}

func ProjectUserProtect(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	projectID := c.Params("projectID")

	var project models.Project
	if err := initializers.DB.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	if project.UserID.String() != userID {
		return &fiber.Error{Code: 403, Message: "Not Allowed to Perform this Action."}
	}

	return c.Next()
}
