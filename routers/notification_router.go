package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func NotificationRouter(app *fiber.App) {
	notificationRoutes := app.Group("/notifications", middlewares.Protect)
	notificationRoutes.Get("/", controllers.GetNotifications)

	notificationRoutes.Get("/unread/count", controllers.GetUnreadNotificationCount)

	notificationRoutes.Get("/unread", controllers.GetUnreadNotifications)

	notificationRoutes.Delete("/:notificationID", controllers.DeleteNotification)
}
