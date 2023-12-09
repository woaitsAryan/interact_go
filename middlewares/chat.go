package middlewares

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func GroupChatAuthorization(Role models.GroupChatRole) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		groupChatID := c.Params("groupChatID")

		var chatMembership models.GroupChatMembership
		if err := initializers.DB.
			First(&chatMembership, "group_chat_id = ? AND user_id = ?", groupChatID, loggedInUserID).Error; err != nil {
			return &fiber.Error{Code: 400, Message: "No chat of this id found."}
		}

		if Role == models.ChatAdmin && chatMembership.Role == models.ChatMember {
			return &fiber.Error{Code: 403, Message: "You do not have the permission to perform this action."}
		}

		return c.Next()
	}
}
