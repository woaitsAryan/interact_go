package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func OpeningRouter(app *fiber.App) {
	openingRoutes := app.Group("/openings", middlewares.Protect)
	openingRoutes.Get("/:openingID", controllers.GetOpening)
	openingRoutes.Get("/project/:projectID", controllers.GetAllOpeningsOfProject)

	openingRoutes.Post("/:projectID", controllers.AddOpening)
	openingRoutes.Patch("/:projectID", controllers.EditOpening)
	openingRoutes.Delete("/:openingID", controllers.DeleteOpening)

}
