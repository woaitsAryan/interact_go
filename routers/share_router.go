package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ShareRouter(app *fiber.App) {

	shareRoutes := app.Group("/share", middlewares.Protect)
	shareRoutes.Post("/post", controllers.SharePost)
	shareRoutes.Post("/project", controllers.ShareProject)

}
