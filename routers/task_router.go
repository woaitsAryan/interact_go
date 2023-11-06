package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func TaskRouter(app *fiber.App) {

	taskRoutes := app.Group("/task", middlewares.Protect)
	taskRoutes.Get("/:taskID", middlewares.ProjectRoleAuthorization(models.ProjectMember), controllers.GetTask("task"))
	taskRoutes.Post("/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.AddTask("task"))
	taskRoutes.Patch("/:taskID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.EditTask("task"))
	taskRoutes.Delete("/:taskID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.DeleteTask("task"))

	taskRoutes.Patch("/completed/:taskID", controllers.MarkTaskCompleted("task")) //* Access Check inside controller
	taskRoutes.Post("/users/:taskID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.AddTaskUser("task"))
	taskRoutes.Delete("/users/:taskID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.RemoveTaskUser("task"))

	taskRoutes.Post("/sub/:taskID", middlewares.TaskUsersAuthorization, controllers.AddTask("subtask"))
	taskRoutes.Patch("/sub/:taskID", middlewares.TaskUsersAuthorization, controllers.EditTask("subtask"))
	taskRoutes.Delete("/sub/:taskID", middlewares.TaskUsersAuthorization, controllers.DeleteTask("subtask"))

	taskRoutes.Patch("/sub/completed/:taskID", controllers.MarkTaskCompleted("subtask")) //* Access Check inside controller
	taskRoutes.Post("/sub/users/:taskID", middlewares.TaskUsersAuthorization, controllers.AddTaskUser("subtask"))
	taskRoutes.Delete("/sub/users/:taskID", middlewares.TaskUsersAuthorization, controllers.RemoveTaskUser("subtask"))
}
