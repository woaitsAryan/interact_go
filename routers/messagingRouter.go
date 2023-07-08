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
	messagingRoutes.Post("/project", controllers.AddProjectChat)

	messagingRoutes.Patch("/group/:chatID", controllers.EditGroupChat)
	messagingRoutes.Patch("/project/:projectChatID", controllers.EditProjectChat)

	messagingRoutes.Delete("/:chatID", controllers.DeleteChat)
	messagingRoutes.Delete("/group/:chatID", controllers.DeleteGroupChat)
	messagingRoutes.Delete("/project/:projectChatID", controllers.DeleteProjectChat)

	messagingRoutes.Get("/content/:chatID", controllers.GetMessages)
	messagingRoutes.Get("/group/content/:projectChatID", controllers.GetMessages)
	messagingRoutes.Get("/project/content/:projectChatID", controllers.GetMessages)

	messagingRoutes.Post("/content", controllers.AddMessage)
	messagingRoutes.Post("/project/content", controllers.AddMessage)

	messagingRoutes.Delete("/content/:messageID", controllers.DeleteMessage)
	messagingRoutes.Delete("/project/content/:messageID", controllers.DeleteMessage)
}
