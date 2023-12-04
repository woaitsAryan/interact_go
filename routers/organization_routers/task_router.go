package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func TaskRouter(app *fiber.App) {

	taskRoutes := app.Group("/org/tasks", middlewares.Protect)
	taskRoutes.Get("/:taskID", middlewares.OrgRoleAuthorization(models.Member), controllers.GetTask("task"))
	taskRoutes.Post("/:orgID", middlewares.OrgRoleAuthorization(models.Manager), controllers.AddTask("org_task"))
	taskRoutes.Patch("/:taskID", middlewares.OrgRoleAuthorization(models.Manager), controllers.EditTask("task"))
	taskRoutes.Delete("/:taskID", middlewares.OrgRoleAuthorization(models.Manager), controllers.DeleteTask("task"))

	taskRoutes.Patch("/completed/:taskID", controllers.MarkTaskCompleted("task")) //* Access Check inside controller
	taskRoutes.Patch("/users/:taskID", middlewares.OrgRoleAuthorization(models.Manager), controllers.AddTaskUser("task"))
	taskRoutes.Delete("/users/:taskID/:userID", middlewares.OrgRoleAuthorization(models.Manager), controllers.RemoveTaskUser("task"))

	taskRoutes.Post("/sub/:taskID", middlewares.TaskUsersCheck, controllers.AddTask("subtask"))
	taskRoutes.Patch("/sub/:taskID", middlewares.SubTaskUsersAuthorization, controllers.EditTask("subtask"))
	taskRoutes.Delete("/sub/:taskID", middlewares.SubTaskUsersAuthorization, controllers.DeleteTask("subtask"))

	taskRoutes.Patch("/sub/completed/:taskID", controllers.MarkTaskCompleted("subtask")) //* Access Check inside controller
	taskRoutes.Patch("/sub/users/:taskID", middlewares.SubTaskUsersAuthorization, controllers.AddTaskUser("subtask"))
	taskRoutes.Delete("/sub/users/:taskID/:userID", middlewares.SubTaskUsersAuthorization, controllers.RemoveTaskUser("subtask"))
}
