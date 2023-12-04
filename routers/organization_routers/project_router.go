package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {

	projectRoutes := app.Group("/org/:orgID/projects", middlewares.Protect)
	projectRoutes.Post("/", middlewares.OrgRoleAuthorization(models.Manager), controllers.AddProject)

	projectRoutes.Patch("/:projectID", middlewares.OrgRoleAuthorization(models.Senior), controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", middlewares.OrgRoleAuthorization(models.Manager), controllers.DeleteProject)
}
