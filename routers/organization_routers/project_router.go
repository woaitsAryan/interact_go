package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {

	projectRoutes := app.Group("/org/projects", middlewares.Protect, middlewares.RoleAuthorization(models.Manager))
	projectRoutes.Post("/", controllers.AddProject)

	projectRoutes.Patch("/:projectID", controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", controllers.DeleteProject)
}
