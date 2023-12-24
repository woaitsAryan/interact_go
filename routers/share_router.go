package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/messaging_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ShareRouter(app *fiber.App) {
	shareRoutes := app.Group("/share", middlewares.Protect)
	shareRoutes.Post("/post", messaging_controllers.ShareItem("post"))
	shareRoutes.Post("/project", messaging_controllers.ShareItem("project"))
	shareRoutes.Post("/opening", messaging_controllers.ShareItem("opening"))
	shareRoutes.Post("/profile", messaging_controllers.ShareItem("profile"))
	shareRoutes.Post("/event", messaging_controllers.ShareItem("event"))
}
