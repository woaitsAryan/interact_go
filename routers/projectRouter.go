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
	projectRoutes.Get("/me/likes", controllers.GetMyLikedProjects)
	// can just have where user_id = logged_in user while searching for project instead of having user-project middleware

	projectRoutes.Get("/:projectID", controllers.GetWorkSpaceProject) //! Add project member protect
	projectRoutes.Patch("/:projectID", controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", middlewares.ProjectUserProtect, controllers.DeleteProject)
	projectRoutes.Get("/like/:projectID", controllers.LikeProject)
}
