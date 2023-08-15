package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ProjectOpeningRouter(app *fiber.App) {

	app.Get("/openings/:openingID", controllers.GetOpening)

	openingRoutes := app.Group("/org/openings", middlewares.Protect, middlewares.RoleAuthorization("Manager"))
	openingRoutes.Get("/project/:projectID", controllers.GetAllOpeningsOfProject)
	openingRoutes.Get("/applications/:openingID", controllers.GetAllApplicationsOfOpening)
	openingRoutes.Post("/:projectID", controllers.AddOpening)
	openingRoutes.Patch("/:openingID", controllers.EditOpening)
	openingRoutes.Delete("/:openingID", controllers.DeleteOpening)
}
