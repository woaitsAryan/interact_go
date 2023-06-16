package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {
	projectRoutes := app.Group("/projects", middlewares.Protect)
	projectRoutes.Post("/", controllers.AddProject)
	projectRoutes.Get("/me", controllers.GetMyProjects)
	projectRoutes.Get("/:projectID", controllers.GetProject)
	projectRoutes.Patch("/:projectID", middlewares.ProjectUserProtect, controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", middlewares.ProjectUserProtect, controllers.DeleteProject)

	projectRoutes.Get("/like/:projectID", controllers.LikeProject)
}
