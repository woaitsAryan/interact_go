package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ShareRouter(app *fiber.App) {

	shareRoutes := app.Group("/share", middlewares.Protect)
	shareRoutes.Post("/post", controllers.ShareItem("post"))
	shareRoutes.Post("/project", controllers.ShareItem("project"))
	shareRoutes.Post("/opening", controllers.ShareItem("opening"))
	shareRoutes.Post("/profile", controllers.ShareItem("profile"))
}
