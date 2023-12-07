package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func WorkspaceRouter(app *fiber.App) {
	workspaceRoutes := app.Group("/workspace", middlewares.Protect)
	workspaceRoutes.Get("/my", project_controllers.GetMyProjects)
	workspaceRoutes.Get("/contributing", project_controllers.GetMyContributingProjects)
	workspaceRoutes.Get("/applications", project_controllers.GetMyApplications)
	workspaceRoutes.Get("/memberships", project_controllers.GetMyMemberships)
}
