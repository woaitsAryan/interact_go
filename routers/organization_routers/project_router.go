package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {
	projectRoutes := app.Group("/org/:orgID/projects", middlewares.Protect)
	projectRoutes.Post("/", middlewares.OrgRoleAuthorization(models.Manager), project_controllers.AddProject)

	projectRoutes.Get("/tasks/:slug", middlewares.OrgRoleAuthorization(models.Senior), project_controllers.GetWorkSpaceProjectTasks)
	projectRoutes.Get("/tasks/populated/:slug", middlewares.OrgRoleAuthorization(models.Senior), project_controllers.GetWorkSpacePopulatedProjectTasks)
	projectRoutes.Get("/history/:projectID", middlewares.OrgRoleAuthorization(models.Member), project_controllers.GetProjectHistory)

	projectRoutes.Get("/:slug", middlewares.OrgRoleAuthorization(models.Member), project_controllers.GetWorkSpaceProject)
	projectRoutes.Patch("/:slug", middlewares.OrgRoleAuthorization(models.Senior), project_controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", middlewares.OrgRoleAuthorization(models.Manager), project_controllers.DeleteProject)
	projectRoutes.Get("/delete/:projectID", middlewares.OrgRoleAuthorization(models.Manager), project_controllers.SendDeleteVerificationCode)
}
