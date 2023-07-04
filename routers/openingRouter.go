package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func OpeningRouter(app *fiber.App) {

	app.Get("/openings/:openingID", controllers.GetOpening)
	app.Get("/openings/project/:projectID", controllers.GetAllOpeningsOfProject)

	openingRoutes := app.Group("/openings", middlewares.Protect)
	openingRoutes.Post("/:projectID", controllers.AddOpening)
	openingRoutes.Patch("/:projectID", controllers.EditOpening)
	openingRoutes.Delete("/:openingID", controllers.DeleteOpening)
}
