package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MessagingRouter(app *fiber.App) {
	messagingRoutes := app.Group("/messaging", middlewares.Protect)

	messagingRoutes.Get("/me", controllers.GetUserNonPopulatedChats)

	messagingRoutes.Get("/", controllers.GetUserChats)
	messagingRoutes.Get("/:chatID", controllers.GetChat)
	messagingRoutes.Get("/project/:projectChatID", controllers.GetChat)

	messagingRoutes.Get("/accept/:chatID", controllers.AcceptChat)

	messagingRoutes.Post("/chat", controllers.AddChat)
	messagingRoutes.Post("/group", controllers.AddGroupChat)
	messagingRoutes.Post("/project/:projectID", controllers.AddProjectChat)

	messagingRoutes.Patch("/group/:chatID", controllers.EditGroupChat)
	messagingRoutes.Patch("/project/:projectChatID", controllers.EditProjectChat)

	messagingRoutes.Delete("/:chatID", controllers.DeleteChat)
	messagingRoutes.Delete("/group/:chatID", controllers.DeleteGroupChat)
	messagingRoutes.Delete("/project/:projectChatID", controllers.DeleteProjectChat)

	messagingRoutes.Get("/content/:chatID", controllers.GetMessages)
	messagingRoutes.Get("/content/group/:ChatID", controllers.GetMessages)
	messagingRoutes.Get("/content/project/:projectChatID", controllers.GetProjectChatMessages)

	messagingRoutes.Post("/content", controllers.AddMessage)
	messagingRoutes.Post("/content/project", controllers.AddProjectChatMessage)

	messagingRoutes.Delete("/content/:messageID", controllers.DeleteMessage)
	messagingRoutes.Delete("/content/project/:messageID", controllers.DeleteMessage)
}
