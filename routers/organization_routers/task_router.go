package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/controllers/organization_controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func TaskRouter(app *fiber.App) {

	app.Get("/org/:orgID/tasks", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Member), organization_controllers.GetOrganizationTasks)

	taskRoutes := app.Group("/org/:orgID/tasks", middlewares.Protect, middlewares.OrgRoleAuthorization(models.Senior))
	taskRoutes.Get("/:taskID", controllers.GetTask("task"))
	taskRoutes.Post("/", controllers.AddTask("org_task"))
	taskRoutes.Patch("/:taskID", controllers.EditTask("task"))
	taskRoutes.Delete("/:taskID", controllers.DeleteTask("task"))

	taskRoutes.Patch("/users/:taskID", controllers.AddTaskUser("org_task"))
	taskRoutes.Delete("/users/:taskID/:userID", controllers.RemoveTaskUser("task"))
}
