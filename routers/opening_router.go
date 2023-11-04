package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func OpeningRouter(app *fiber.App) {

	app.Get("/openings/:openingID", middlewares.PartialProtect, controllers.GetOpening)

	openingRoutes := app.Group("/openings", middlewares.Protect)
	openingRoutes.Get("/project/:projectID", controllers.GetAllOpeningsOfProject)
	openingRoutes.Get("/applications/:openingID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.GetAllApplicationsOfOpening)
	openingRoutes.Post("/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.AddOpening)
	openingRoutes.Patch("/:openingID", middlewares.ProjectRoleAuthorization(models.ProjectEditor), controllers.EditOpening)
	openingRoutes.Delete("/:openingID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.DeleteOpening)
}
