package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func WorkspaceRouter(app *fiber.App) {
	workspaceRoutes := app.Group("/workspace", middlewares.Protect)
	workspaceRoutes.Get("/my", controllers.GetMyProjects)
	workspaceRoutes.Get("/contributing", controllers.GetMyContributingProjects)
	workspaceRoutes.Get("/applications", controllers.GetMyApplications)
	workspaceRoutes.Get("/memberships", controllers.GetMyMemberships)
}
