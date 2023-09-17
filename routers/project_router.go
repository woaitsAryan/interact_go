package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {

	projectRoutes := app.Group("/projects", middlewares.Protect)
	projectRoutes.Post("/", controllers.AddProject)
	projectRoutes.Get("/me", controllers.GetMyProjects)
	projectRoutes.Get("/me/likes", controllers.GetMyLikedProjects)
	projectRoutes.Get("/:slug", middlewares.ProjectRoleAuthorization(models.ProjectMember), controllers.GetWorkSpaceProject)
	projectRoutes.Get("/chats/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectEditor), controllers.GetWorkSpaceProjectChats)
	projectRoutes.Patch("/:slug", middlewares.ProjectRoleAuthorization(models.ProjectEditor), controllers.UpdateProject)
	projectRoutes.Delete("/:projectID", controllers.DeleteProject)
	projectRoutes.Get("/like/:projectID", controllers.LikeProject)
}
