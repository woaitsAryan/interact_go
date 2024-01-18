package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func FeedRouter(app *fiber.App) {
	feedRoutes := app.Group("/feed", middlewares.Protect)
	feedRoutes.Get("/", controllers.GetFeed)
	feedRoutes.Get("/combined", controllers.GetCombinedFeed)
}
