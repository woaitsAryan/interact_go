package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {

	projectRoutes := app.Group("/org/projects", middlewares.Protect, middlewares.RoleAuthorization("Manager"))
	projectRoutes.Post("/", controllers.AddProject)
	projectRoutes.Get("/me", controllers.GetMyProjects)
	projectRoutes.Get("/me/likes", controllers.GetMyLikedProjects)

	projectRoutes.Get("/:projectID", controllers.GetWorkSpaceProject)
	projectRoutes.Patch("/:projectID", controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", controllers.DeleteProject)
	projectRoutes.Get("/like/:projectID", controllers.LikeProject)
}
