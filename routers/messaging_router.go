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

	messagingRoutes.Get("/", controllers.GetUserChats)
	messagingRoutes.Get("/personal", controllers.GetUserChats)
	messagingRoutes.Get("/:chatID", controllers.GetChat)
	messagingRoutes.Get("/group/:groupChat", controllers.GetGroupChat)

	messagingRoutes.Get("/accept/:chatID", controllers.AcceptChat)

	messagingRoutes.Post("/chat", controllers.AddChat)
	messagingRoutes.Post("/group", controllers.AddGroupChat)
	messagingRoutes.Post("/project/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.AddProjectChat)

	messagingRoutes.Patch("/group/:chatID", controllers.EditGroupChat)

	messagingRoutes.Delete("/:chatID", controllers.DeleteChat)

	messagingRoutes.Delete("/group/:chatID", controllers.DeleteGroupChat)

	messagingRoutes.Delete("/group/:chatID", controllers.LeaveGroupChat)

	messagingRoutes.Get("/content/:chatID", controllers.GetMessages)
	messagingRoutes.Get("/content/group/:ChatID", controllers.GetGroupChatMessages)

	messagingRoutes.Post("/content", controllers.AddMessage)
	messagingRoutes.Post("/content/group", controllers.AddGroupChatMessage)

	messagingRoutes.Delete("/content/:messageID", controllers.DeleteMessage)
	messagingRoutes.Delete("/content/project/:messageID", controllers.DeleteMessage)
}
