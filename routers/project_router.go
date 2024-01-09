package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/project_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func ProjectRouter(app *fiber.App) {
	projectRoutes := app.Group("/projects", middlewares.Protect)
	projectRoutes.Post("/", project_controllers.AddProject)
	projectRoutes.Get("/me", project_controllers.GetMyProjects)
	projectRoutes.Get("/me/likes", project_controllers.GetMyLikedProjects)
	projectRoutes.Get("/:slug", middlewares.ProjectRoleAuthorization(models.ProjectMember), project_controllers.GetWorkSpaceProject)
	projectRoutes.Get("/chats/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectEditor), project_controllers.GetWorkSpaceProjectChats)
	projectRoutes.Patch("/:slug", middlewares.ProjectRoleAuthorization(models.ProjectEditor), project_controllers.UpdateProject)
	projectRoutes.Get("/like/:projectID", controllers.LikeItem("project"))

	projectRoutes.Get("/history/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectMember), project_controllers.GetProjectHistory)
	projectRoutes.Get("/tasks/:slug", middlewares.ProjectRoleAuthorization(models.ProjectMember), project_controllers.GetWorkSpaceProjectTasks)
	projectRoutes.Get("/tasks/populated/:slug", middlewares.ProjectRoleAuthorization(models.ProjectMember), project_controllers.GetWorkSpacePopulatedProjectTasks)

	projectRoutes.Get("/delete/:projectID", project_controllers.SendDeleteVerificationCode)
	projectRoutes.Delete("/:projectID", project_controllers.DeleteProject)
}
