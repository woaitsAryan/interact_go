package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/messaging_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

//TODO34 have separate routers and controllers for project group n organization chats

func MessagingRouter(app *fiber.App) {
	messagingRoutes := app.Group("/messaging", middlewares.Protect)

	messagingRoutes.Get("/me", messaging_controllers.GetUserNonPopulatedChats)

	messagingRoutes.Get("/personal", messaging_controllers.GetPersonalChats)
	messagingRoutes.Get("/personal/unfiltered", messaging_controllers.GetPersonalUnFilteredChats)
	messagingRoutes.Get("/personal/unread", messaging_controllers.GetUnreadChats)
	messagingRoutes.Get("/group", messaging_controllers.GetGroupChats)
	messagingRoutes.Get("/project", messaging_controllers.GetProjectChats)

	messagingRoutes.Get("/:chatID", messaging_controllers.GetChat)
	messagingRoutes.Get("/group/:chatID", messaging_controllers.GetGroupChat)

	messagingRoutes.Get("/accept/:chatID", messaging_controllers.AcceptChat)

	messagingRoutes.Post("/chat", messaging_controllers.AddChat)
	messagingRoutes.Post("/group", messaging_controllers.AddGroupChat("Group"))
	messagingRoutes.Post("/project/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectEditor), messaging_controllers.AddGroupChat("Project"))

	messagingRoutes.Patch("/chat/last_read/:chatID", messaging_controllers.UpdateLastRead)

	messagingRoutes.Post("/chat/block", messaging_controllers.BlockChat)
	messagingRoutes.Post("/chat/unblock", messaging_controllers.UnblockChat)
	messagingRoutes.Post("/chat/reset", messaging_controllers.ResetChat)

	messagingRoutes.Post("/group/members/add/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.AddGroupChatMembers("Group"))
	messagingRoutes.Post("/group/members/remove/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.RemoveGroupChatMember)

	messagingRoutes.Patch("/group/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.EditGroupChat)
	messagingRoutes.Patch("/group/role/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.EditGroupChatRole)

	messagingRoutes.Delete("/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.DeleteChat)
	messagingRoutes.Delete("/group/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.DeleteGroupChat)

	messagingRoutes.Delete("/group/leave/:chatID", messaging_controllers.LeaveGroupChat) //TODO35 when admin leaves, then make the first joined person as admin

	messagingRoutes.Get("/content/:chatID", messaging_controllers.GetMessages)
	messagingRoutes.Get("/content/group/:chatID", messaging_controllers.GetGroupChatMessages)

	messagingRoutes.Post("/content", messaging_controllers.AddMessage)
	messagingRoutes.Post("/content/group", messaging_controllers.AddGroupChatMessage)

	messagingRoutes.Delete("/content/:messageID", messaging_controllers.DeleteMessage)
	messagingRoutes.Delete("/content/project/:messageID", messaging_controllers.DeleteMessage)

	messagingRoutes.Post("/group/project/members/add/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.AddGroupChatMembers("Project"))
	messagingRoutes.Post("/group/project/members/remove/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.RemoveGroupChatMember)

	messagingRoutes.Post("/group/organization/members/add/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.AddGroupChatMembers("Organization"))
	messagingRoutes.Post("/group/organization/members/remove/:chatID", middlewares.GroupChatAdminAuthorization(), messaging_controllers.RemoveGroupChatMember)
}
