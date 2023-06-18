package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MessagingRouter(app *fiber.App) {
	messagingRoutes := app.Group("/messaging", middlewares.Protect)

	messagingRoutes.Get("/:chatID", controllers.GetChat)
	messagingRoutes.Get("/project/:projectChatID", controllers.GetChat)

	messagingRoutes.Post("/chat", controllers.AddChat)
	messagingRoutes.Post("/group", controllers.AddGroupChat)
	messagingRoutes.Post("/project", controllers.AddProjectChat)

	messagingRoutes.Patch("/group/:chatID", controllers.EditGroupChat)
	messagingRoutes.Patch("/project/:projectChatID", controllers.EditProjectChat)

	messagingRoutes.Delete("/:chatID", controllers.DeleteChat)
	messagingRoutes.Delete("/project/:projectChatID", controllers.DeleteProjectChat)

	messagingRoutes.Get("/content/:chatID", controllers.GetMessages)
	messagingRoutes.Post("/content", controllers.AddMessage)
	messagingRoutes.Delete("/content/:messageID", controllers.DeleteMessage)
}
