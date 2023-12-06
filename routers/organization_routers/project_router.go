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

	projectRoutes.Get("/tasks/:slug", middlewares.OrgRoleAuthorization(models.Senior), controllers.GetWorkSpaceProjectTasks)
	projectRoutes.Get("/tasks/populated/:slug", middlewares.OrgRoleAuthorization(models.Senior), controllers.GetWorkSpacePopulatedProjectTasks)
	projectRoutes.Get("/history/:projectID", middlewares.OrgRoleAuthorization(models.Member), controllers.GetProjectHistory)

	projectRoutes.Get("/:slug", middlewares.OrgRoleAuthorization(models.Member), controllers.GetWorkSpaceProject)
	projectRoutes.Patch("/:slug", middlewares.OrgRoleAuthorization(models.Senior), controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", middlewares.OrgRoleAuthorization(models.Manager), controllers.DeleteProject)
}
