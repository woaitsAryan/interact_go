package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

//TODO have separate routers and controllers for project group n organization chats

func MessagingRouter(app *fiber.App) {
	messagingRoutes := app.Group("/messaging", middlewares.Protect)

	messagingRoutes.Get("/me", controllers.GetUserNonPopulatedChats)

	messagingRoutes.Get("/personal", controllers.GetPersonalChats)
	messagingRoutes.Get("/personal/unfiltered", controllers.GetPersonalUnFilteredChats)
	messagingRoutes.Get("/group", controllers.GetGroupChats)
	messagingRoutes.Get("/project", controllers.GetProjectChats)

	messagingRoutes.Get("/:chatID", controllers.GetChat)
	messagingRoutes.Get("/group/:chatID", controllers.GetGroupChat)

	messagingRoutes.Get("/accept/:chatID", controllers.AcceptChat)

	messagingRoutes.Post("/chat", controllers.AddChat)
	messagingRoutes.Post("/group", controllers.AddGroupChat)
	messagingRoutes.Post("/project/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectEditor), controllers.AddProjectChat)

	messagingRoutes.Patch("/chat/last_read/:chatID", controllers.UpdateLastRead)

	messagingRoutes.Post("/chat/block", controllers.BlockChat)
	messagingRoutes.Post("/chat/unblock", controllers.UnblockChat)
	messagingRoutes.Post("/chat/reset", controllers.ResetChat)

	messagingRoutes.Post("/group/members/add/:chatID", middlewares.GroupChatAdminAuthorization(), controllers.AddGroupChatMembers)
	messagingRoutes.Post("/group/members/remove/:chatID", middlewares.GroupChatAdminAuthorization(), controllers.RemoveGroupChatMember)

	messagingRoutes.Patch("/group/:chatID", middlewares.GroupChatAdminAuthorization(), controllers.EditGroupChat)
	messagingRoutes.Patch("/group/role/:chatID", middlewares.GroupChatAdminAuthorization(), controllers.EditGroupChatRole)

	messagingRoutes.Delete("/:chatID", controllers.DeleteChat)
	messagingRoutes.Delete("/group/:chatID", middlewares.GroupChatAdminAuthorization(), controllers.DeleteGroupChat)

	messagingRoutes.Delete("/group/leave/:chatID", controllers.LeaveGroupChat)

	messagingRoutes.Get("/content/:chatID", controllers.GetMessages)
	messagingRoutes.Get("/content/group/:chatID", controllers.GetGroupChatMessages)

	messagingRoutes.Post("/content", controllers.AddMessage)
	messagingRoutes.Post("/content/group", controllers.AddGroupChatMessage)

	messagingRoutes.Delete("/content/:messageID", controllers.DeleteMessage)
	messagingRoutes.Delete("/content/project/:messageID", controllers.DeleteMessage)

	messagingRoutes.Post("/group/project/members/add/:chatID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.AddProjectChatMembers)
	messagingRoutes.Post("/group/project/members/remove/:chatID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.RemoveGroupChatMember)

	messagingRoutes.Patch("/group/project/:chatID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.EditGroupChat)
	messagingRoutes.Patch("/group/project/role/:chatID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.EditGroupChatRole)

	messagingRoutes.Delete("/group/project/:chatID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.DeleteGroupChat)

}
