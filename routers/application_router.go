package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ApplicationRouter(app *fiber.App) {
	applicationRoutes := app.Group("/applications", middlewares.Protect)
	applicationRoutes.Post("/:openingID", controllers.AddApplication)

	applicationRoutes.Get("/accept/:applicationID", controllers.AcceptApplication)
	applicationRoutes.Get("/reject/:applicationID", controllers.RejectApplication)
	applicationRoutes.Get("/review/:applicationID", controllers.SetApplicationUnderReview)

	applicationRoutes.Get("/:applicationID", controllers.GetApplication)
	applicationRoutes.Delete("/:applicationID", controllers.DeleteApplication)
}
