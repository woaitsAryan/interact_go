package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ProjectOpeningRouter(app *fiber.App) {

	app.Get("/openings/:openingID", project_controllers.GetOpening)

	openingRoutes := app.Group("/org/:orgID/openings", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	openingRoutes.Get("/applications/:openingID", controllers.GetAllApplicationsOfOpening)
	openingRoutes.Post("/:projectID", project_controllers.AddOpening)
	openingRoutes.Patch("/:openingID", project_controllers.EditOpening)
	openingRoutes.Delete("/:openingID", project_controllers.DeleteOpening)
}
